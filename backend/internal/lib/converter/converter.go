package converter

import (
	"time"

	"github.com/mickamy/automapper"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	userv1 "github.com/mickamy/sampay/gen/user/v1"
)

func init() {
	automapper.RegisterFromE[userv1.PaymentMethodType, string](ToPaymentMethodType)
	automapper.RegisterFrom[string, *string](StringToPtr)
	automapper.RegisterFrom[int32, int](Int32ToInt)
	automapper.RegisterFrom[int, int32](IntToInt32)
	automapper.RegisterFrom[time.Time, *timestamppb.Timestamp](TimeToTimestamppb)
	automapper.RegisterFrom[string, eventv1.ParticipantStatus](ToV1ParticipantStatus)
	automapper.RegisterFrom[eventv1.ParticipantStatus, string](FromV1ParticipantStatus)
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

func ToV1ParticipantStatus(s string) eventv1.ParticipantStatus {
	switch s {
	case "unpaid":
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNPAID
	case "claimed":
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_CLAIMED
	case "confirmed":
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED
	default:
		return eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNSPECIFIED
	}
}

func FromV1ParticipantStatus(s eventv1.ParticipantStatus) string {
	switch s {
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNPAID:
		return "unpaid"
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_CLAIMED:
		return "claimed"
	case eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED:
		return "confirmed"
	default:
		return ""
	}
}
