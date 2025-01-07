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
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestAccount_SignUp(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.SignUpRequest
		assert  func(t *testing.T, got *connect.Response[registrationv1.SignUpResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.SignUpRequest {
				return &registrationv1.SignUpRequest{
					Email:    gofakeit.GlobalFaker.Email(),
					Password: gofakeit.GlobalFaker.Password(true, true, true, true, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.SignUpResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.UserId)
				assert.NotEmpty(t, got.Msg.Tokens.Access.Value)
				assert.NotEmpty(t, got.Msg.Tokens.Access.ExpiresAt)
				assert.NotEmpty(t, got.Msg.Tokens.Refresh.Value)
				assert.NotEmpty(t, got.Msg.Tokens.Refresh.ExpiresAt)
			},
		},
		{
			name: "email already exists",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *registrationv1.SignUpRequest {
				user := userFixture.User(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = user.ID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				return &registrationv1.SignUpRequest{
					Email:    auth.Identifier,
					Password: gofakeit.GlobalFaker.Password(true, true, true, true, false, 12),
				}
			},
			assert: func(t *testing.T, got *connect.Response[registrationv1.SignUpResponse], err error) {
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
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: "registration.handler.error.email_already_exists"}), errMsg.FieldViolations[0].Descriptions[0])
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
			server := newAccountServer(t, infras)

			// act
			client := registrationv1connect.NewAccountServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewRequest(t, ctx, req, nil)
			got, err := client.SignUp(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newAccountServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitRegistrationHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).Account
		return registrationv1connect.NewAccountServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
