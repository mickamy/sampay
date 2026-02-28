package converter

import (
	"errors"
	"fmt"
	"time"

	"github.com/mickamy/automapper"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/domain/event/model"
)

func init() {
	automapper.RegisterFromE[userv1.PaymentMethodType, string](ToPaymentMethodType)
	automapper.RegisterFrom[string, *string](StringToPtr)
	automapper.RegisterFrom[int32, int](Int32ToInt)
	automapper.RegisterFrom[int, int32](IntToInt32)
	automapper.RegisterFrom[time.Time, *timestamppb.Timestamp](TimeToTimestamppb)
	automapper.RegisterFrom[*time.Time, *timestamppb.Timestamp](PtrTimeToTimestamppb)
	automapper.RegisterFrom[model.ParticipantStatus, eventv1.ParticipantStatus](ToV1ParticipantStatus)
	automapper.RegisterFromE[eventv1.ParticipantStatus, model.ParticipantStatus](FromV1ParticipantStatus)
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

func IntToInt32(i int) int32 {
	return int32(i) //nolint:gosec // values are reasonable small positive integers
}

func TimeToTimestamppb(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func PtrTimeToTimestamppb(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func ToV1ParticipantStatus(s model.ParticipantStatus) eventv1.ParticipantStatus {
	switch s {
	case model.ParticipantStatusUnpaid:
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNPAID
	case model.ParticipantStatusClaimed:
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_CLAIMED
	case model.ParticipantStatusConfirmed:
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED
	default:
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNSPECIFIED
	}
}

func FromV1ParticipantStatus(s eventv1.ParticipantStatus) (model.ParticipantStatus, error) {
	switch s {
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNPAID:
		return model.ParticipantStatusUnpaid, nil
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_CLAIMED:
		return model.ParticipantStatusClaimed, nil
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED:
		return model.ParticipantStatusConfirmed, nil
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNSPECIFIED:
		return "", errors.New("unspecified participant status")
	default:
		return "", fmt.Errorf("unknown participant status: %v", s)
	}
}
