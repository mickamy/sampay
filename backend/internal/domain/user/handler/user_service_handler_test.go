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

func TestUserService_GetMe(t *testing.T) {
	t.Parallel()

	t.Run("returns current user", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.GetMeResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceGetMeProcedure).
			Header("Authorization", authHeader).
			In(&userv1.GetMeRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.NotEmpty(t, out.GetUser().GetId())
		assert.NotEmpty(t, out.GetUser().GetSlug())
	})
}

func TestUserService_UpdateSlug(t *testing.T) {
	t.Parallel()

	t.Run("updates slug successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.UpdateSlugResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceUpdateSlugProcedure).
			Header("Authorization", authHeader).
			In(&userv1.UpdateSlugRequest{Slug: "my-new-slug"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, "my-new-slug", out.GetUser().GetSlug())
	})

	t.Run("returns error for invalid slug", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceUpdateSlugProcedure).
			Header("Authorization", authHeader).
			In(&userv1.UpdateSlugRequest{Slug: "AB"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
	})

	t.Run("returns error when slug already taken", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) {
			m.UserID = user.ID
			m.Slug = "taken-slug"
		})
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceUpdateSlugProcedure).
			Header("Authorization", authHeader).
			In(&userv1.UpdateSlugRequest{Slug: "taken-slug"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusConflict)
	})
}

func TestUserService_CheckSlugAvailability(t *testing.T) {
	t.Parallel()

	t.Run("returns available when slug is free", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.CheckSlugAvailabilityResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceCheckSlugAvailabilityProcedure).
			Header("Authorization", authHeader).
			In(&userv1.CheckSlugAvailabilityRequest{Slug: "free-slug"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.True(t, out.GetAvailable())
	})

	t.Run("returns unavailable when slug is taken", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) {
			m.UserID = user.ID
			m.Slug = "taken-slug"
		})
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		_, authHeader := ctest.UserSession(t, infra)

		// act
		var out userv1.CheckSlugAvailabilityResponse
		ct := contest.NewWith(t,
			contest.Bind(userv1connect.NewUserServiceHandler)(handler.NewUserService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(userv1connect.UserServiceCheckSlugAvailabilityProcedure).
			Header("Authorization", authHeader).
			In(&userv1.CheckSlugAvailabilityRequest{Slug: "taken-slug"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.False(t, out.GetAvailable())
	})
}
