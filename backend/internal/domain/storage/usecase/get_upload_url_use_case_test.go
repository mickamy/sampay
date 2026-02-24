package usecase_test

import (
	"context"
	"testing"

	"github.com/mickamy/enufstub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/domain/storage/usecase"
	"github.com/mickamy/sampay/internal/infra/aws/s3"
)

func TestGetUploadURL_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns presigned URL and creates S3Object record", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := enufstub.Of[s3.Client]().With("PresignPutObject",
			func(_ context.Context, bucket, key string) (string, error) {
				assert.Equal(t, "s3-public-bucket", bucket)
				assert.Equal(t, "qr/user1/paypay.png", key)
				return "https://s3.example.com/presigned", nil
			}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})

		// act
		sut := usecase.NewGetUploadURL(infra)
		out, err := sut.Do(t.Context(), usecase.GetUploadURLInput{
			Path: "qr/user1/paypay.png",
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, "https://s3.example.com/presigned", out.UploadURL)
		assert.NotEmpty(t, out.S3ObjectID)

		// verify DB record
		obj, err := query.S3Objects(infra.ReaderDB).Where("id = ?", out.S3ObjectID).First(t.Context())
		require.NoError(t, err)
		assert.Equal(t, "qr/user1/paypay.png", obj.Key)
	})

	t.Run("returns error when S3 presign fails", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := enufstub.Of[s3.Client]().With("PresignPutObject",
			func(_ context.Context, _, _ string) (string, error) {
				return "", assert.AnError
			}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})

		// act
		sut := usecase.NewGetUploadURL(infra)
		_, err := sut.Do(t.Context(), usecase.GetUploadURLInput{
			Path: "qr/user1/paypay.png",
		})

		// assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to generate presigned URL")
	})
}
