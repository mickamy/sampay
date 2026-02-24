package mapper

import (
	"github.com/mickamy/errx"

	userv1 "github.com/mickamy/sampay/gen/user/v1"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

const (
	PaymentMethodTypePayPay     = "paypay"
	PaymentMethodTypeKyash      = "kyash"
	PaymentMethodTypeRakutenPay = "rakuten_pay"
	PaymentMethodTypeMerPay     = "merpay"
)

func ToPaymentMethodType(src userv1.PaymentMethodType) (string, error) {
	switch src {
	case userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED:
		// fall through to error
	case userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY:
		return PaymentMethodTypePayPay, nil
	case userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH:
		return PaymentMethodTypeKyash, nil
	case userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_RAKUTEN_PAY:
		return PaymentMethodTypeRakutenPay, nil
	case userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY:
		return PaymentMethodTypeMerPay, nil
	}
	return "", cmodel.NewLocalizableError(
		errx.New("unknown payment method type", "type", src).WithCode(errx.InvalidArgument),
	).
		WithMessages(messages.UserMapperErrorUnknownPaymentMethodType())
}

func ToV1PaymentMethodType(src string) userv1.PaymentMethodType {
	switch src {
	case PaymentMethodTypePayPay:
		return userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY
	case PaymentMethodTypeKyash:
		return userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH
	case PaymentMethodTypeRakutenPay:
		return userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_RAKUTEN_PAY
	case PaymentMethodTypeMerPay:
		return userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY
	default:
		return userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED
	}
}
