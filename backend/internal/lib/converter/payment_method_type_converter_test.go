package converter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/lib/converter"
)

func TestToPaymentMethodType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     userv1.PaymentMethodType
		want    string
		wantErr bool
	}{
		{"paypay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY, "paypay", false},
		{"kyash", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH, "kyash", false},
		{"rakuten_pay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_RAKUTEN_PAY, "rakuten_pay", false},
		{"merpay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY, "merpay", false},
		{"unspecified", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := converter.ToPaymentMethodType(tt.src)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestToV1PaymentMethodType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		src  string
		want userv1.PaymentMethodType
	}{
		{"paypay", "paypay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_PAYPAY},
		{"kyash", "kyash", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_KYASH},
		{"rakuten_pay", "rakuten_pay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_RAKUTEN_PAY},
		{"merpay", "merpay", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_MERPAY},
		{"unknown", "unknown", userv1.PaymentMethodType_PAYMENT_METHOD_TYPE_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := converter.ToV1PaymentMethodType(tt.src)
			assert.Equal(t, tt.want, got)
		})
	}
}
