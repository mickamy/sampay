package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/domain/storage/usecase"
)

type mockS3Client struct {
	url string
	err error
}

func (m *mockS3Client) PresignPutObject(_ context.Context, _, _ string) (string, error) {
	return m.url, m.err
}

func TestGetUploadURL_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns presigned URL and creates S3Object record", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := &mockS3Client{url: "https://s3.example.com/presigned"}
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock
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
		mock := &mockS3Client{err: assert.AnError}
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock
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
