package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestPasswordReset_ResetPassword(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest
		assert  func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest {
				verification := authFixture.EmailVerificationVerified(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
					m.Identifier = verification.Email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)

				return &authv1.ResetPasswordRequest{
					Token:       verification.Verified.Token,
					NewPassword: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "fail (invalid token)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest {
				verification := authFixture.EmailVerificationVerified(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
					m.Identifier = verification.Email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)

				return &authv1.ResetPasswordRequest{
					Token:       verification.Verified.Token + "invalid",
					NewPassword: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error) {
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
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseReset_passwordErrorEmail_verification_invalid_token}), errMsg.FieldViolations[0].Descriptions[0])
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
		{
			name: "fail (token of request)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest {
				verification := authFixture.EmailVerificationVerified(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
					m.Identifier = verification.Email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)

				return &authv1.ResetPasswordRequest{
					Token:       verification.Requested.Token,
					NewPassword: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error) {
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
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseReset_passwordErrorEmail_verification_invalid_token}), errMsg.FieldViolations[0].Descriptions[0])
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
		{
			name: "fail (token consumed)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest {
				verification := authFixture.EmailVerificationConsumed(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
					m.Identifier = verification.Email
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)

				return &authv1.ResetPasswordRequest{
					Token:       verification.Verified.Token,
					NewPassword: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error) {
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
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseReset_passwordErrorEmail_verification_already_consumed}), errMsg.FieldViolations[0].Descriptions[0])
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
		{
			name: "fail (no authentication)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.ResetPasswordRequest {
				verification := authFixture.EmailVerificationVerified(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&verification).Error)

				return &authv1.ResetPasswordRequest{
					Token:       verification.Verified.Token,
					NewPassword: gofakeit.Password(true, true, true, false, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.ResetPasswordResponse], err error) {
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
			server := newResetPasswordServer(t, infras)

			// act
			client := authv1connect.NewPasswordResetServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.ResetPassword(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newResetPasswordServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitAuthHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).PasswordReset
		return authv1connect.NewPasswordResetServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
