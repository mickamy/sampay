package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/common/v1/commonv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *commonv1.CreateDirectUploadURLRequest
		assert  func(t *testing.T, got *connect.Response[commonv1.CreateDirectUploadURLResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *commonv1.CreateDirectUploadURLRequest {
				return &commonv1.CreateDirectUploadURLRequest{
					S3Object: &commonv1.S3Object{
						Bucket: gofakeit.GlobalFaker.ProductName(),
						Key:    gofakeit.UUID(),
					},
				}
			},
			assert: func(t *testing.T, got *connect.Response[commonv1.CreateDirectUploadURLResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.Url)
			},
		},
		{
			name: "fail",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *commonv1.CreateDirectUploadURLRequest {
				return &commonv1.CreateDirectUploadURLRequest{
					S3Object: nil,
				}
			},
			assert: func(t *testing.T, got *connect.Response[commonv1.CreateDirectUploadURLResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.CommonHandlerDirect_upload_urlErrorInvalid_s3_object})
					assert.Equal(t, expectedMsg, errMsg.Message)
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
			server := newSessionServer(t, infras)
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			req := tc.arrange(t, ctx, infras)

			// act
			client := commonv1connect.NewDirectUploadURLServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.CreateDirectUploadURL(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newSessionServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitCommonHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).DirectUploadURL
		return commonv1connect.NewDirectUploadURLServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
