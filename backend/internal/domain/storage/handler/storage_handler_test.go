package handler_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/enufstub"
	"github.com/mickamy/errx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/mickamy/sampay/gen/storage/v1"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/handler"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/infra/aws/s3"
)

func TestStorage_GetUploadURL(t *testing.T) {
	t.Parallel()

	t.Run("returns presigned URL and s3 object ID", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := enufstub.Of[s3.Client]().With("PresignPutObject", func(_ context.Context, bucket, key string) (string, error) {
			assert.Equal(t, "sampay-public", bucket)
			assert.Equal(t, "qr/user1/paypay.png", key)
			return "https://s3.example.com/presigned", nil
		}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})
		sut := handler.NewStorage(infra)

		// act
		res, err := sut.GetUploadURL(t.Context(), connect.NewRequest(&v1.GetUploadURLRequest{
			Path: "qr/user1/paypay.png",
		}))

		// assert
		require.NoError(t, err)
		assert.Equal(t, "https://s3.example.com/presigned", res.Msg.GetUploadUrl())
		assert.NotEmpty(t, res.Msg.GetS3ObjectId())

		// verify DB record
		obj, err := query.S3Objects(infra.ReaderDB).Where("id = ?", res.Msg.GetS3ObjectId()).First(t.Context())
		require.NoError(t, err)
		assert.Equal(t, "qr/user1/paypay.png", obj.Key)
	})

	t.Run("returns error when S3 presign fails", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := enufstub.Of[s3.Client]().With("PresignPutObject", func(_ context.Context, _, _ string) (string, error) {
			return "", assert.AnError
		}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})
		sut := handler.NewStorage(infra)

		// act
		_, err := sut.GetUploadURL(t.Context(), connect.NewRequest(&v1.GetUploadURLRequest{
			Path: "qr/user1/paypay.png",
		}))

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.Internal))
	})
}
