package scenario

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/lib/ptr"
)

func TestOnboarding(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	infras := di.NewInfras(newReadWriter(t), newKVS(t))
	server := initServer(t, infras)

	token := emailVerification(t, ctx, infras, server)
	onboarding(t, ctx, infras, server, token)
	assert.NotEmpty(t, token)
}

func getOnboardingStep(t *testing.T, s *httptest.Server, verifyToken string, f func(res *connect.Response[registrationv1.GetOnboardingStepResponse], err error)) {
	t.Helper()

	client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, s.URL)
	req := connect.NewRequest(&registrationv1.GetOnboardingStepRequest{})
	req.Header().Add("Authorization", "Bearer "+verifyToken)
	res, err := client.GetOnboardingStep(context.Background(), req)
	f(res, err)
}

func createPassword(t *testing.T, s *httptest.Server, verifyToken string, f func(res *connect.Response[registrationv1.CreatePasswordResponse], err error)) {
	t.Helper()

	client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, s.URL)
	req := connect.NewRequest(&registrationv1.CreatePasswordRequest{
		Password: password,
	})
	req.Header().Add("Authorization", "Bearer "+verifyToken)
	res, err := client.CreatePassword(context.Background(), req)
	f(res, err)
}

func getUsageCategories(t *testing.T, s *httptest.Server, accessToken string, f func(res *connect.Response[registrationv1.ListUsageCategoriesResponse], err error)) {
	t.Helper()

	client := registrationv1connect.NewUsageCategoryServiceClient(http.DefaultClient, s.URL)
	req := connect.NewRequest(&registrationv1.ListUsageCategoriesRequest{})
	req.Header().Add("Authorization", "Bearer "+accessToken)
	res, err := client.ListUsageCategories(context.Background(), req)
	f(res, err)
}

func createUserAttribute(t *testing.T, s *httptest.Server, accessToken string, f func(res *connect.Response[registrationv1.CreateUserAttributeResponse], err error)) {
	t.Helper()

	client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, s.URL)
	req := connect.NewRequest(&registrationv1.CreateUserAttributeRequest{
		CategoryType: "other",
	})
	req.Header().Add("Authorization", "Bearer "+accessToken)
	res, err := client.CreateUserAttribute(context.Background(), req)
	f(res, err)
}

func createUserProfile(t *testing.T, s *httptest.Server, accessToken string, f func(res *connect.Response[registrationv1.CreateUserProfileResponse], err error)) {
	t.Helper()

	client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, s.URL)
	req := connect.NewRequest(&registrationv1.CreateUserProfileRequest{
		Name: gofakeit.GlobalFaker.Name(),
		Bio:  ptr.Of(gofakeit.GlobalFaker.Sentence(10)),
		Image: &commonv1.S3Object{
			Bucket: gofakeit.GlobalFaker.ProductName(),
			Key:    gofakeit.GlobalFaker.UUID(),
		},
	})
	req.Header().Add("Authorization", "Bearer "+accessToken)
	res, err := client.CreateUserProfile(context.Background(), req)
	f(res, err)
}

func onboarding(t *testing.T, ctx context.Context, infras di.Infras, s *httptest.Server, verifyToken string) string {
	t.Helper()

	{
		getOnboardingStep(t, s, verifyToken, func(res *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
			require.NoError(t, err)
			assert.Equal(t, "password", res.Msg.Step)
		})
	}

	var accessToken string
	{
		createPassword(t, s, verifyToken, func(res *connect.Response[registrationv1.CreatePasswordResponse], err error) {
			require.NoError(t, err)
			require.NotEmpty(t, res.Msg.Tokens)
			accessToken = res.Msg.Tokens.Access.Value
		})
	}

	{
		getOnboardingStep(t, s, verifyToken, func(res *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
			require.NoError(t, err)
			assert.Equal(t, "attribute", res.Msg.Step)
		})
	}

	{
		getUsageCategories(t, s, accessToken, func(res *connect.Response[registrationv1.ListUsageCategoriesResponse], err error) {
			require.NoError(t, err)
			require.NotEmpty(t, res.Msg.Categories)
		})
	}

	{
		createUserAttribute(t, s, accessToken, func(res *connect.Response[registrationv1.CreateUserAttributeResponse], err error) {
			require.NoError(t, err)
		})
	}

	{
		getOnboardingStep(t, s, verifyToken, func(res *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
			require.NoError(t, err)
			assert.Equal(t, "profile", res.Msg.Step)
		})
	}

	{
		createUserProfile(t, s, accessToken, func(res *connect.Response[registrationv1.CreateUserProfileResponse], err error) {
			require.NoError(t, err)
		})
	}

	{
		getOnboardingStep(t, s, verifyToken, func(res *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
			require.NoError(t, err)
			assert.Equal(t, "completed", res.Msg.Step)
		})
	}

	return accessToken
}
