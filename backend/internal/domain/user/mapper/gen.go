package mapper

import (
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/domain/user/model"
)

//go:generate go tool automapper -from=model.EndUser -to=userv1.User -output=./ -converter-pkg=../../../lib/converter
var (
	_ model.EndUser
	_ userv1.User
)
