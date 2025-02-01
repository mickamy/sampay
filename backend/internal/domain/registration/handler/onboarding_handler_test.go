package handler_test

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
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/lib/random"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestOnboarding_GetOnboardingStep(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, email string) *registrationv1.GetOnboardingStepRequest
		assert  func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error)
	}{
		{
			name: "success (password)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, email string) *registrationv1.GetOnboardingStepRequest {
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepPassword.String(), got.Msg.Step)
			},
		},
		{
			name: "success (attribute)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, email string) *registrationv1.GetOnboardingStepRequest {
				user := userFixture.User(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = user.ID
					m.Identifier = email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepAttribute.String(), got.Msg.Step)
			},
		},
		{
			name: "success (profile)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, email string) *registrationv1.GetOnboardingStepRequest {
				user := userFixture.User(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = user.ID
					m.Identifier = email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				attr := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = user.ID
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
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, email string) *registrationv1.GetOnboardingStepRequest {
				user := userFixture.User(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = user.ID
					m.Identifier = email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				attr := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = user.ID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&attr).Error)
				profile := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = user.ID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&profile).Error)
				return &registrationv1.GetOnboardingStepRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
				require.NoError(t, err)
				require.Equal(t, registrationModel.OnboardingStepCompleted.String(), got.Msg.Step)
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
			verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
				m.IntentType = authModel.EmailVerificationIntentTypeSignUp
			})
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)
			req := tc.arrange(t, ctx, infras, verification.Email)
			server := newOnboardingServer(t, infras)

			// act
			client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAnonymousRequest(t, ctx, req, nil, verification)
			got, err := client.GetOnboardingStep(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestOnboarding_CreateUserPassword(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))

	tsc := []struct {
		name         string
		arrange      func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.CreatePasswordRequest
		authenticate func(req *connect.Request[registrationv1.CreatePasswordRequest])
		assert       func(t *testing.T, got *connect.Response[registrationv1.CreatePasswordResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.CreatePasswordRequest {
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)

				return &registrationv1.CreatePasswordRequest{
					Password: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			authenticate: func(req *connect.Request[registrationv1.CreatePasswordRequest]) {
				req.Header().Add("Authorization", "Bearer "+token)
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreatePasswordResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.Tokens)
			},
		},
		{
			name: "fail (invalid token)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.CreatePasswordRequest {
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)

				return &registrationv1.CreatePasswordRequest{
					Password: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			authenticate: func(req *connect.Request[registrationv1.CreatePasswordRequest]) {
				req.Header().Add("Authorization", "Bearer "+token+"invalid")
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreatePasswordResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeUnauthenticated, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 0)
			},
		},
		{
			name: "fail (token consumed)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.CreatePasswordRequest {
				verification := authFixture.EmailVerificationConsumed(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)

				return &registrationv1.CreatePasswordRequest{
					Password: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			authenticate: func(req *connect.Request[registrationv1.CreatePasswordRequest]) {
				req.Header().Add("Authorization", "Bearer "+token)
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreatePasswordResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
					require.Len(t, errMsg.FieldViolations, 1)
					require.Equal(t, "token", errMsg.FieldViolations[0].Field)
					require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_passwordErrorEmail_verification_already_consumed}), errMsg.FieldViolations[0].Descriptions[0])
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
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
			req := tc.arrange(t, ctx, infras)
			server := newOnboardingServer(t, infras)

			// act
			client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewRequest(t, ctx, req, nil)
			tc.authenticate(connReq)
			got, err := client.CreatePassword(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestOnboarding_CreateUserAttribute(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserAttributeRequest
		assert  func(t *testing.T, got *connect.Response[registrationv1.CreateUserAttributeResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserAttributeRequest {
				return &registrationv1.CreateUserAttributeRequest{
					CategoryType: "other",
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreateUserAttributeResponse], err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "fail (invalid category type)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserAttributeRequest {
				return &registrationv1.CreateUserAttributeRequest{
					CategoryType: "invalid",
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreateUserAttributeResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInternal, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.CommonHandlerErrorInternal}), errMsg.Message)
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
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
			got, err := client.CreateUserAttribute(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestOnboarding_CreateUserProfile(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserProfileRequest
		assert  func(t *testing.T, got *connect.Response[registrationv1.CreateUserProfileResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserProfileRequest {
				return &registrationv1.CreateUserProfileRequest{
					Name: gofakeit.GlobalFaker.Name(),
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreateUserProfileResponse], err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "success with bio",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserProfileRequest {
				return &registrationv1.CreateUserProfileRequest{
					Name: gofakeit.GlobalFaker.Name(),
					Bio:  ptr.Of(gofakeit.GlobalFaker.Sentence(20)),
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreateUserProfileResponse], err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "success with image",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *registrationv1.CreateUserProfileRequest {
				s3Obj := commonFixture.S3Object(nil)
				return &registrationv1.CreateUserProfileRequest{
					Name: gofakeit.GlobalFaker.Name(),
					Image: &commonv1.S3Object{
						Bucket:      s3Obj.Bucket,
						Key:         s3Obj.Key,
						ContentType: s3Obj.ContentType.String(),
					},
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.CreateUserProfileResponse], err error) {
				require.NoError(t, err)
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
			got, err := client.CreateUserProfile(ctx, connReq)

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
