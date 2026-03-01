package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mickamy/sampay/internal/di"
	amodel "github.com/mickamy/sampay/internal/domain/auth/model"
	arepository "github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/domain/messaging/lib/line"
	mrepository "github.com/mickamy/sampay/internal/domain/messaging/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

type ClaimNotificationPayload struct {
	CreatorUserID   string `json:"creator_user_id" validate:"required"`
	EventTitle      string `json:"event_title"     validate:"required"`
	ParticipantName string `json:"participant_name" validate:"required"`
	Amount          int    `json:"amount"           validate:"required"`
}

type ClaimNotification struct {
	_                  *di.Infra                  `inject:"param"`
	reader             *database.Reader           `inject:""`
	oauthAccountRepo   arepository.OAuthAccount   `inject:""`
	lineFriendshipRepo mrepository.LineFriendship `inject:""`
	pushClient         *line.PushClient           `inject:""`
}

func (j *ClaimNotification) Execute(ctx context.Context, payloadStr string) error {
	var payload ClaimNotificationPayload
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	var uid string
	var isFriend bool

	if err := j.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		uid, err = j.oauthAccountRepo.WithTx(tx).GetUIDByEndUserIDAndProvider(ctx, payload.CreatorUserID, amodel.OAuthProviderLINE)
		if err != nil {
			return fmt.Errorf("failed to get LINE UID for user %s: %w", payload.CreatorUserID, err)
		}

		friendship, err := j.lineFriendshipRepo.WithTx(tx).GetByEndUserID(ctx, payload.CreatorUserID)
		if err != nil {
			return fmt.Errorf("failed to get friendship for user %s: %w", payload.CreatorUserID, err)
		}
		isFriend = friendship.IsFriend
		return nil
	}); err != nil {
		return err //nolint:wrapcheck // already wrapped inside
	}

	if !isFriend {
		logger.Info(ctx, "user is not a LINE friend; skipping push", "user_id", payload.CreatorUserID)
		return nil
	}

	text := i18n.Japanese(messages.MessagingClaimNotification(payload.ParticipantName, payload.EventTitle, payload.Amount))

	if err := j.pushClient.SendTextMessage(ctx, uid, text); err != nil {
		return fmt.Errorf("failed to send LINE push message: %w", err)
	}

	return nil
}
