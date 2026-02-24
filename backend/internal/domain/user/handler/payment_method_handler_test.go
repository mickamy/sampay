package handler_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/contest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/gen/user/v1/userv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/handler"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/test/ctest"
)

func TestPaymentMethod_ListPaymentMethods(t *testing.T) {
	t.Parallel()

	t.Run("returns payment methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		m := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = userID })
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &m))

		// act
		var out userv1.ListPaymentMethodsResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceListPaymentMethodsProcedure).
			Header("Authorization", authHeader).
			In(&userv1.ListPaymentMethodsRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Len(t, out.GetPaymentMethods(), 1)
		assert.Equal(t, m.ID, out.GetPaymentMethods()[0].GetId())
	})

	t.Run("returns empty when no methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.ListPaymentMethodsResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceListPaymentMethodsProcedure).
			Header("Authorization", authHeader).
			In(&userv1.ListPaymentMethodsRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Empty(t, out.GetPaymentMethods())
	})

	t.Run("returns unauthenticated without token", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceListPaymentMethodsProcedure).
			In(&userv1.ListPaymentMethodsRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusUnauthorized)
	})
}

func TestPaymentMethod_SavePaymentMethods(t *testing.T) {
	t.Parallel()

	t.Run("saves payment methods and returns them", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.SavePaymentMethodsResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceSavePaymentMethodsProcedure).
			Header("Authorization", authHeader).
			In(&userv1.SavePaymentMethodsRequest{
				PaymentMethods: []*userv1.PaymentMethodInput{
					{
						Type:         userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY,
						Url:          "https://paypay.example.com",
						DisplayOrder: 0,
					},
					{
						Type:         userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH,
						Url:          "https://kyash.example.com",
						DisplayOrder: 1,
					},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Len(t, out.GetPaymentMethods(), 2)
		assert.Equal(t, userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY, out.GetPaymentMethods()[0].GetType())
		assert.Equal(t, userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH, out.GetPaymentMethods()[1].GetType())
	})

	t.Run("returns error for invalid payment method type", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceSavePaymentMethodsProcedure).
			Header("Authorization", authHeader).
			In(&userv1.SavePaymentMethodsRequest{
				PaymentMethods: []*userv1.PaymentMethodInput{
					{
						Type:         userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED,
						Url:          "https://example.com",
						DisplayOrder: 0,
					},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
	})

	t.Run("replaces existing payment methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)

		existing := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = userID })
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &existing))

		// act
		var out userv1.SavePaymentMethodsResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewPaymentMethodServiceHandler)(handler.NewPaymentMethod(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.PaymentMethodServiceSavePaymentMethodsProcedure).
			Header("Authorization", authHeader).
			In(&userv1.SavePaymentMethodsRequest{
				PaymentMethods: []*userv1.PaymentMethodInput{
					{
						Type:         userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY,
						Url:          "https://merpay.example.com",
						DisplayOrder: 0,
					},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Len(t, out.GetPaymentMethods(), 1)
		assert.Equal(t, userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY, out.GetPaymentMethods()[0].GetType())
	})
}
