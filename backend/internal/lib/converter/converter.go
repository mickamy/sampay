package converter

import (
	"github.com/mickamy/automapper"

	userv1 "github.com/mickamy/sampay/gen/user/v1"
)

func init() {
	automapper.RegisterFromE[string, userv1.PaymentMethodType](ToPaymentMethodType)
	automapper.RegisterFrom[*string, string](StringToPtr)
	automapper.RegisterFrom[int, int32](Int32ToInt)
}

func StringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Int32ToInt(i int32) int {
	return int(i)
}
