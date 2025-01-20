package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestEmailVerification_RequestVerification(t *testing.T) {
	t.Parallel()

	t.Run("sign_up", func(t *testing.T) {
		t.Parallel()

		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.RequestVerificationRequest
			assert  func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error)
		}{
			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.RequestVerificationRequest {
					return &authv1.RequestVerificationRequest{
						Email: gofakeit.GlobalFaker.Email(),
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Msg.Token)
				},
			},
			{
				name: "email already exists",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.RequestVerificationRequest {
					user := userFixture.User(nil)
					require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
					auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
						m.UserID = user.ID
					})
					require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
					return &authv1.RequestVerificationRequest{
						Email: auth.Identifier,
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
						require.Len(t, errMsg.FieldViolations, 1)
						require.Equal(t, "email", errMsg.FieldViolations[0].Field)
						require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
						require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_accountErrorEmail_already_exists}), errMsg.FieldViolations[0].Descriptions[0])
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
				req.IntentType = authv1.RequestVerificationRequest_INTENT_TYPE_SIGN_UP
				server := newEmailVerificationServer(t, infras)

				// act
				client := authv1connect.NewEmailVerificationServiceClient(http.DefaultClient, server.URL)
				connReq := connecttest.NewRequest(t, ctx, req, nil)
				got, err := client.RequestVerification(ctx, connReq)

				// assert
				tc.assert(t, got, err)
			})
		}
	})

	t.Run("reset_password", func(t *testing.T) {
		t.Parallel()

		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RequestVerificationRequest
			assert  func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error)
		}{
			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RequestVerificationRequest {
					auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
						m.UserID = userID
					})
					require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
					return &authv1.RequestVerificationRequest{
						Email: auth.Identifier,
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Msg.Token)
				},
			},
			{
				name: "no authentication",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RequestVerificationRequest {
					return &authv1.RequestVerificationRequest{
						Email: gofakeit.GlobalFaker.Email(),
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RequestVerificationResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
						require.Len(t, errMsg.FieldViolations, 1)
						require.Equal(t, "email", errMsg.FieldViolations[0].Field)
						require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
						require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseRequest_email_verificationErrorEmail_not_found}), errMsg.FieldViolations[0].Descriptions[0])
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
				ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
				req := tc.arrange(t, ctx, infras, user.ID)
				req.IntentType = authv1.RequestVerificationRequest_INTENT_TYPE_RESET_PASSWORD
				server := newEmailVerificationServer(t, infras)

				// act
				client := authv1connect.NewEmailVerificationServiceClient(http.DefaultClient, server.URL)
				connReq := connecttest.NewRequest(t, ctx, req, nil)
				got, err := client.RequestVerification(ctx, connReq)

				// assert
				tc.assert(t, got, err)
			})
		}
	})
}

func TestEmailVerification_VerifyEmail(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.VerifyEmailRequest
		assert  func(t *testing.T, got *connect.Response[authv1.VerifyEmailResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.VerifyEmailRequest {
				request := authFixture.EmailVerificationRequested(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&request).Error)
				return &authv1.VerifyEmailRequest{
					Token:   request.Requested.Token,
					PinCode: request.Requested.PINCode,
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.VerifyEmailResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.Session)
				assert.NotEmpty(t, got.Msg.Token)
			},
		},
		{
			name: "invalid token",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *authv1.VerifyEmailRequest {
				request := authFixture.EmailVerificationRequested(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&request).Error)
				return &authv1.VerifyEmailRequest{
					Token:   request.Requested.Token,
					PinCode: either.Must(random.NewPinCode(6)),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.VerifyEmailResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
					require.Len(t, errMsg.FieldViolations, 1)
					require.Equal(t, "pin_code", errMsg.FieldViolations[0].Field)
					require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseVerify_emailErrorInvalid_pin_code}), errMsg.FieldViolations[0].Descriptions[0])
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
			server := newEmailVerificationServer(t, infras)

			// act
			client := authv1connect.NewEmailVerificationServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewRequest(t, ctx, req, nil)
			got, err := client.VerifyEmail(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newEmailVerificationServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitAuthHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).EmailVerification
		return authv1connect.NewEmailVerificationServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
