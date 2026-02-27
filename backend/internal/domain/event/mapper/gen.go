package mapper

import (
	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/internal/domain/event/model"
)

//go:generate go tool automapper -from=model.Event -to=eventv1.Event -output=./ -converter-pkg=../../../lib/converter
//go:generate go tool automapper -from=model.EventParticipant -to=eventv1.EventParticipant -output=./ -converter-pkg=../../../lib/converter
var (
	_ model.Event
	_ eventv1.Event
	_ model.EventParticipant
	_ eventv1.EventParticipant
)
