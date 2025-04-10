package scenario_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	userModel "mickamy.com/sampay/internal/domain/user/model"
)

func TestUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	infras := di.NewInfras(newReadWriter(t), newKVS(t))
	server := initServer(t, infras)

	verificationToken := emailVerification(t, ctx, infras, server)
	accessToken := onboarding(t, ctx, infras, server, verificationToken)
	getMe(t, server, accessToken, func(res *connect.Response[userv1.GetMeResponse], err error) {
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Msg.User)
		assert.NotEmpty(t, res.Msg.User.Profile)
		assert.Len(t, res.Msg.User.Links, 0)
	})

	var slug string
	{
		infras.DB.Model(&userModel.User{}).Order("id DESC").Limit(1).Pluck("slug", &slug)
	}

	getUser(t, server, slug, func(res *connect.Response[userv1.GetUserResponse], err error) {
		require.NoError(t, err)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Msg.User)
		assert.NotEmpty(t, res.Msg.User.Profile)
		assert.Len(t, res.Msg.User.Links, 0)
	})
}

func getMe(t *testing.T, s *httptest.Server, accessToken string, f func(res *connect.Response[userv1.GetMeResponse], err error)) {
	t.Helper()

	client := userv1connect.NewUserServiceClient(http.DefaultClient, s.URL+"/api")
	req := connect.NewRequest(&userv1.GetMeRequest{})
	req.Header().Add("Authorization", "Bearer "+accessToken)
	res, err := client.GetMe(context.Background(), req)
	f(res, err)
}

func getUser(t *testing.T, s *httptest.Server, slug string, f func(res *connect.Response[userv1.GetUserResponse], err error)) {
	t.Helper()

	client := userv1connect.NewUserServiceClient(http.DefaultClient, s.URL+"/api")
	req := connect.NewRequest(&userv1.GetUserRequest{
		Slug: slug,
	})
	res, err := client.GetUser(context.Background(), req)
	f(res, err)
}
