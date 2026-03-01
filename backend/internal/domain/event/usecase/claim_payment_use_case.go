package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"
	"github.com/mickamy/go-sqs-worker/message"
	"github.com/mickamy/go-sqs-worker/producer"

	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/job"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrClaimPaymentNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("participant not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorParticipantNotFound())
	ErrClaimPaymentAlreadyClaimed = cmodel.NewLocalizableError(
		errx.NewSentinel("already claimed", errx.FailedPrecondition),
	).WithMessages(messages.EventUseCaseErrorAlreadyClaimed())
	ErrClaimPaymentArchived = cmodel.NewLocalizableError(
		errx.NewSentinel("event is archived", errx.FailedPrecondition),
	).WithMessages(messages.EventUseCaseErrorArchived())
)

type ClaimPaymentInput struct {
	ParticipantID string
}

type ClaimPaymentOutput struct {
	Participant model.EventParticipant
}

type ClaimPayment interface {
	Do(ctx context.Context, input ClaimPaymentInput) (ClaimPaymentOutput, error)
}

type claimPayment struct {
	_               ClaimPayment                `inject:"returns"`
	_               *di.Infra                   `inject:"param"`
	writer          *database.Writer            `inject:""`
	eventRepo       repository.Event            `inject:""`
	participantRepo repository.EventParticipant `inject:""`
	producer        *producer.Producer          `inject:""`
}

func (uc *claimPayment) Do(ctx context.Context, input ClaimPaymentInput) (ClaimPaymentOutput, error) {
	var participant model.EventParticipant
	var ev model.Event

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		var err error
		participant, err = uc.participantRepo.WithTx(tx).Get(ctx, input.ParticipantID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrClaimPaymentNotFound
			}
			return errx.Wrap(err, "message", "failed to get participant", "id", input.ParticipantID).
				WithCode(errx.Internal)
		}

		var evErr error
		ev, evErr = uc.eventRepo.WithTx(tx).Get(ctx, participant.EventID)
		if evErr != nil {
			return errx.Wrap(evErr, "message", "failed to get event", "event_id", participant.EventID).
				WithCode(errx.Internal)
		}
		if ev.ArchivedAt != nil {
			return ErrClaimPaymentArchived
		}

		if participant.Status != model.ParticipantStatusUnpaid {
			return ErrClaimPaymentAlreadyClaimed
		}

		participant.Status = model.ParticipantStatusClaimed
		if err := uc.participantRepo.WithTx(tx).Update(ctx, &participant); err != nil {
			return errx.Wrap(err, "message", "failed to update participant", "id", input.ParticipantID).
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return ClaimPaymentOutput{}, err
	}

	uc.enqueueClaimNotification(ctx, ev, participant)

	return ClaimPaymentOutput{Participant: participant}, nil
}

func (uc *claimPayment) enqueueClaimNotification(
	ctx context.Context,
	ev model.Event,
	participant model.EventParticipant,
) {
	msg, err := message.New(ctx, job.ClaimNotificationJob, job.ClaimNotificationPayload{
		CreatorUserID:   ev.UserID,
		EventTitle:      ev.Title,
		ParticipantName: participant.Name,
		Amount:          participant.Amount,
	})
	if err != nil {
		logger.Error(ctx, "failed to create claim notification message",
			"error", err, "event_id", ev.ID, "participant_id", participant.ID)
		return
	}

	if err := uc.producer.Do(ctx, msg); err != nil {
		logger.Error(ctx, "failed to enqueue claim notification",
			"error", err, "event_id", ev.ID, "participant_id", participant.ID)
	}
}
