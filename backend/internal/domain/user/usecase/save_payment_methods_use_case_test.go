package usecase_test

import (
	"testing"

	"github.com/mickamy/errx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sfixture "github.com/mickamy/sampay/internal/domain/storage/fixture"
	smodel "github.com/mickamy/sampay/internal/domain/storage/model"
	squery "github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

func TestSavePaymentMethods_Do(t *testing.T) {
	t.Parallel()

	t.Run("creates new payment methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		out, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "https://pay.example.com/paypay", DisplayOrder: 0},
				{Type: "kyash", URL: "https://pay.example.com/kyash", DisplayOrder: 1},
			},
		})

		// assert
		require.NoError(t, err)
		assert.Len(t, out.PaymentMethods, 2)
	})

	t.Run("replaces existing payment methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		existing := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "paypay"
		})
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &existing))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		out, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "kyash", URL: "https://pay.example.com/kyash", DisplayOrder: 0},
			},
		})

		// assert
		require.NoError(t, err)
		assert.Len(t, out.PaymentMethods, 1)
		assert.Equal(t, "kyash", out.PaymentMethods[0].Type)
	})

	t.Run("rejects duplicate types", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		_, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "https://a.com", DisplayOrder: 0},
				{Type: "paypay", URL: "https://b.com", DisplayOrder: 1},
			},
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "duplicate payment method type")
	})

	t.Run("rejects empty URL", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		_, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "", DisplayOrder: 0},
			},
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "payment method URL is required")
	})

	t.Run("rejects empty type", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		_, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "", URL: "https://a.com", DisplayOrder: 0},
			},
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "payment method type is required")
	})

	t.Run("empty array deletes all existing methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		existing := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "paypay"
		})
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &existing))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		out, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{},
		})

		// assert
		require.NoError(t, err)
		assert.Empty(t, out.PaymentMethods)
	})

	t.Run("accepts own S3 object", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		s3Obj := sfixture.S3Object(func(m *smodel.S3Object) {
			m.Key = user.ID + "/qr/paypay.png"
		})
		require.NoError(t, squery.S3Objects(infra.WriterDB).Create(t.Context(), &s3Obj))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)
		s3ObjID := s3Obj.ID

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		out, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "https://pay.example.com/paypay", QRCodeS3ObjectID: &s3ObjID, DisplayOrder: 0},
			},
		})

		// assert
		require.NoError(t, err)
		assert.Len(t, out.PaymentMethods, 1)
		require.NotNil(t, out.PaymentMethods[0].QRCodeS3ObjectID)
		assert.Equal(t, s3ObjID, *out.PaymentMethods[0].QRCodeS3ObjectID)
	})

	t.Run("rejects S3 object belonging to another user", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		otherUser := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &otherUser))

		s3Obj := sfixture.S3Object(func(m *smodel.S3Object) {
			m.Key = otherUser.ID + "/qr/paypay.png"
		})
		require.NoError(t, squery.S3Objects(infra.WriterDB).Create(t.Context(), &s3Obj))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)
		s3ObjID := s3Obj.ID

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		_, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "https://pay.example.com/paypay", QRCodeS3ObjectID: &s3ObjID, DisplayOrder: 0},
			},
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "does not belong to the authenticated user")
	})

	t.Run("rejects javascript URL scheme", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewSavePaymentMethods(infra)
		_, err := sut.Do(ctx, usecase.SavePaymentMethodsInput{
			PaymentMethods: []usecase.SavePaymentMethodInput{
				{Type: "paypay", URL: "javascript:alert(1)", DisplayOrder: 0},
			},
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "valid HTTP or HTTPS URL")
	})
}
