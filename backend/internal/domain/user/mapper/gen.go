package mapper

import (
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
)

//go:generate go tool automapper -from=model.EndUser -to=userv1.User -output=./ -converter-pkg=../../../lib/converter
//go:generate go tool automapper -from=userv1.PaymentMethodInput -to=usecase.SavePaymentMethodInput -output=./ -converter-pkg=../../../lib/converter
var (
	_ model.EndUser
	_ userv1.User

	_ userv1.PaymentMethodInput
	_ usecase.SavePaymentMethodInput
)
