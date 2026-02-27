package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/lib/converter"
)

// ToV1Participant converts an EventParticipant model to a proto message.
// amount is a computed value (not stored in DB), so it must be passed explicitly.
func ToV1Participant(src model.EventParticipant, amount int) *eventv1.EventParticipant {
	return &eventv1.EventParticipant{
		Id:        src.ID,
		EventId:   src.EventID,
		Name:      src.Name,
		Tier:      int32(src.Tier),   //nolint:gosec // Tier is 1-5
		Status:    converter.ToV1ParticipantStatus(src.Status),
		Amount:    int32(amount),     //nolint:gosec // Amount is a reasonable positive integer
		CreatedAt: timestamppb.New(src.CreatedAt),
	}
}
