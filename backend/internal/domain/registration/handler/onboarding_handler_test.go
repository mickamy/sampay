package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/test/connecttest"
)

func TestOnboarding_GetOnboardingStep(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.GetOnboardingStepRequest
		assert  func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error)
	}{
		{
			name: "success (attribute)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.GetOnboardingStepRequest {
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepAttribute.String(), got.Msg.Step)
			},
		},
		{
			name: "success (profile)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.GetOnboardingStepRequest {
				attr := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&attr).Error)
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepProfile.String(), got.Msg.Step)
			},
		},
		{
			name: "success (complete)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.GetOnboardingStepRequest {
				attr := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&attr).Error)
				profile := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&profile).Error)
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepComplete.String(), got.Msg.Step)
			},
		},
	}

	for _, tc := range tsc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			infras := di.NewInfras(newReadWriter(t), newKVS(t))
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newOnboardingServer(t, infras)

			// act
			client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.GetOnboardingStep(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newOnboardingServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitRegistrationHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).Onboarding
		return registrationv1connect.NewOnboardingServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
