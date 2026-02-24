package mapper

import (
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/lib/converter"
)

func ToV1PaymentMethod(src model.UserPaymentMethod, cloudfrontBaseURL string) *userv1.PaymentMethod {
	pm := &userv1.PaymentMethod{
		Id:           src.ID,
		Type:         converter.ToV1PaymentMethodType(src.Type),
		Url:          src.URL,
		DisplayOrder: int32(src.DisplayOrder), //nolint:gosec // DisplayOrder is a small non-negative integer
	}
	if src.QRCodeS3ObjectID != nil {
		pm.QrCodeS3ObjectId = *src.QRCodeS3ObjectID
	}
	if src.QRCodeS3Object != nil {
		pm.QrCodeUrl = cloudfrontBaseURL + "/" + src.QRCodeS3Object.Key
	}
	return pm
}
