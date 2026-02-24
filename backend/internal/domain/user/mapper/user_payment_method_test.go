package mapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	userv1 "github.com/mickamy/sampay/gen/user/v1"
	smodel "github.com/mickamy/sampay/internal/domain/storage/model"
	"github.com/mickamy/sampay/internal/domain/user/mapper"
	"github.com/mickamy/sampay/internal/domain/user/model"
)

func TestToV1PaymentMethod(t *testing.T) {
	t.Parallel()

	t.Run("with QR code", func(t *testing.T) {
		t.Parallel()

		src := model.UserPaymentMethod{
			ID:           "id-1",
			Type:         "paypay",
			URL:          "https://example.com/pay",
			DisplayOrder: 1,
			QRCodeS3Object: &smodel.S3Object{
				Key: "qr/user1/paypay.png",
			},
		}

		got := mapper.ToV1PaymentMethod(src, "https://cdn.example.com")

		assert.Equal(t, "id-1", got.Id)
		assert.Equal(t, userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY, got.Type)
		assert.Equal(t, "https://example.com/pay", got.Url)
		assert.Equal(t, int32(1), got.DisplayOrder)
		assert.Equal(t, "https://cdn.example.com/qr/user1/paypay.png", got.QrCodeUrl)
	})

	t.Run("without QR code", func(t *testing.T) {
		t.Parallel()

		src := model.UserPaymentMethod{
			ID:   "id-2",
			Type: "kyash",
			URL:  "https://example.com/kyash",
		}

		got := mapper.ToV1PaymentMethod(src, "https://cdn.example.com")

		assert.Equal(t, "", got.QrCodeUrl)
	})
}

func TestToV1PaymentMethods(t *testing.T) {
	t.Parallel()

	src := []model.UserPaymentMethod{
		{ID: "id-1", Type: "paypay", URL: "https://a.com"},
		{ID: "id-2", Type: "kyash", URL: "https://b.com"},
	}

	got := mapper.ToV1PaymentMethods(src, "https://cdn.example.com")

	assert.Len(t, got, 2)
	assert.Equal(t, "id-1", got[0].Id)
	assert.Equal(t, "id-2", got[1].Id)
}
