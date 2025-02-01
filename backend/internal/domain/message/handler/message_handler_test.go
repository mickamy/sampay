package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/message/v1/messagev1connect"
	messagev1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/message/v1"
	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/test/connecttest"
)

func TestMessage_SendMessage(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *messagev1.SendMessageRequest
		assert  func(t *testing.T, got *connect.Response[messagev1.SendMessageResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *messagev1.SendMessageRequest {
				receiver := userFixture.User(nil)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&receiver).Error)
				return &messagev1.SendMessageRequest{
					SenderName:   gofakeit.GlobalFaker.Name(),
					ReceiverSlug: receiver.Slug,
					Content:      gofakeit.GlobalFaker.Sentence(20),
				}
			},
			assert: func(t *testing.T, got *connect.Response[messagev1.SendMessageResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
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
			server := newMessageServer(t, infras)
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			req := tc.arrange(t, ctx, infras)

			// act
			client := messagev1connect.NewMessageServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.SendMessage(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newMessageServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitMessageHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).Message
		return messagev1connect.NewMessageServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
