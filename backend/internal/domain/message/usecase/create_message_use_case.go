package usecase

import (
	"context"
	"fmt"

	"github.com/mickamy/go-sqs-worker/message"
	"github.com/mickamy/go-sqs-worker/producer"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	messageModel "mickamy.com/sampay/internal/domain/message/model"
	messageRepository "mickamy.com/sampay/internal/domain/message/repository"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/job"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type CreateMessageInput struct {
	SenderName   string
	ReceiverSlug string
	Content      string
}

type CreateMessageOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateMessage interface {
	Do(ctx context.Context, input CreateMessageInput) (CreateMessageOutput, error)
}

type createMessage struct {
	producer    *producer.Producer
	writer      *database.Writer
	userRepo    userRepository.User
	messageRepo messageRepository.Message
}

func NewCreateMessage(
	producer *producer.Producer,
	writer *database.Writer,
	userRepo userRepository.User,
	messageRepo messageRepository.Message,
) CreateMessage {
	return &createMessage{
		producer:    producer,
		writer:      writer,
		userRepo:    userRepo,
		messageRepo: messageRepo,
	}
}

func (uc *createMessage) Do(ctx context.Context, input CreateMessageInput) (CreateMessageOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		receiver, err := uc.userRepo.WithTx(tx.WriterDB()).FindBySlug(ctx, input.ReceiverSlug)
		if err != nil {
			return fmt.Errorf("failed to find receiver: %w", err)
		}
		if receiver == nil {
			return fmt.Errorf("receiver not found: slug=[%s]", input.ReceiverSlug)
		}

		msg := messageModel.Message{
			SenderName: input.SenderName,
			ReceiverID: receiver.ID,
			Content:    input.Content,
		}
		if err := uc.messageRepo.WithTx(tx.WriterDB()).Create(ctx, &msg); err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}

		lang := contexts.MustLanguage(ctx)
		jobMsg, err := message.New(ctx, job.SendEmailJob.String(), job.SendEmailPayload{
			From: config.Email().From,
			To:   receiver.Email,
			Subject: i18n.MustLocalizeMessage(lang, i18n.Config{
				MessageID:    i18n.MessageUsecaseCreate_messageEmailSubject,
				TemplateData: map[string]string{"SenderName": input.SenderName},
			}),
			Body: i18n.MustLocalizeMessage(lang, i18n.Config{
				MessageID:    i18n.MessageUsecaseCreate_messageEmailBody,
				TemplateData: map[string]string{"SenderName": input.SenderName, "Content": input.Content},
			}),
		})
		if err != nil {
			return fmt.Errorf("failed to create job message: %w", err)
		}

		if err := uc.producer.Do(ctx, jobMsg); err != nil {
			return fmt.Errorf("failed to send job message: %w", err)
		}

		return nil
	}); err != nil {
		return CreateMessageOutput{}, err
	}

	return CreateMessageOutput{}, nil
}
