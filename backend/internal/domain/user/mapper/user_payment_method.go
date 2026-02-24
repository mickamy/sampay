package mapper

import (
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/domain/user/model"
)

func ToV1PaymentMethod(src model.UserPaymentMethod, cloudFrontBaseURL string) *userv1.PaymentMethod {
	pm := &userv1.PaymentMethod{
		Id:           src.ID,
		Type:         ToV1PaymentMethodType(src.Type),
		Url:          src.URL,
		DisplayOrder: int32(src.DisplayOrder),
	}
	if src.QRCodeS3Object != nil {
		pm.QrCodeUrl = cloudFrontBaseURL + "/" + src.QRCodeS3Object.Key
	}
	return pm
}

func ToV1PaymentMethods(src []model.UserPaymentMethod, cloudFrontBaseURL string) []*userv1.PaymentMethod {
	result := make([]*userv1.PaymentMethod, len(src))
	for i, m := range src {
		result[i] = ToV1PaymentMethod(m, cloudFrontBaseURL)
	}
	return result
}
