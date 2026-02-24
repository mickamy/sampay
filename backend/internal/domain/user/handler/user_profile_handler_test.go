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
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestUserProfile_GetUserProfile(t *testing.T) {
	t.Parallel()

	t.Run("returns user profile without authentication", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		pm := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = endUser.UserID })
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &pm))

		// act
		var out userv1.GetUserProfileResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserProfileServiceHandler)(handler.NewUserProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserProfileServiceGetUserProfileProcedure).
			In(&userv1.GetUserProfileRequest{Slug: endUser.Slug}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, endUser.UserID, out.GetUser().GetId())
		assert.Equal(t, endUser.Slug, out.GetUser().GetSlug())
		assert.Len(t, out.GetPaymentMethods(), 1)
		assert.Equal(t, pm.ID, out.GetPaymentMethods()[0].GetId())
	})

	t.Run("returns not found for unknown slug", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserProfileServiceHandler)(handler.NewUserProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserProfileServiceGetUserProfileProcedure).
			In(&userv1.GetUserProfileRequest{Slug: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
	})

	t.Run("returns empty payment methods when none registered", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		// act
		var out userv1.GetUserProfileResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserProfileServiceHandler)(handler.NewUserProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserProfileServiceGetUserProfileProcedure).
			In(&userv1.GetUserProfileRequest{Slug: endUser.Slug}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, endUser.Slug, out.GetUser().GetSlug())
		assert.Empty(t, out.GetPaymentMethods())
	})
}
