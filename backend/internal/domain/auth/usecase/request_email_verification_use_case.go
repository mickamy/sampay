package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/go-sqs-worker/message"
	"github.com/mickamy/go-sqs-worker/producer"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/job"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrRequestEmailVerificationEmailAlreadyExists = commonModel.
		NewLocalizableError(errors.New("email already exists")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCommonErrorEmail_already_exists})
)

type RequestEmailVerificationInput struct {
	Email    string
	Password string
}

type RequestEmailVerificationOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type RequestEmailVerification interface {
	Do(ctx context.Context, input RequestEmailVerificationInput) (RequestEmailVerificationOutput, error)
}

type requestEmailVerification struct {
	writer                *database.Writer
	producer              *producer.Producer
	authenticationRepo    authRepository.Authentication
	emailVerificationRepo authRepository.EmailVerification
}

func NewRequestEmailVerification(
	writer *database.Writer,
	producer *producer.Producer,
	authenticationRepo authRepository.Authentication,
	emailVerificationRepo authRepository.EmailVerification,
) RequestEmailVerification {
	return &requestEmailVerification{
		writer:                writer,
		producer:              producer,
		authenticationRepo:    authenticationRepo,
		emailVerificationRepo: emailVerificationRepo,
	}
}

func (uc *requestEmailVerification) Do(ctx context.Context, input RequestEmailVerificationInput) (RequestEmailVerificationOutput, error) {
	m := authModel.EmailVerification{Email: input.Email}
	cfg := config.Auth()
	if err := m.Request(cfg.EmailVerificationExpiresInDuration()); err != nil {
		return RequestEmailVerificationOutput{}, fmt.Errorf("failed to request email verification: %w", err)
	}
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		exists, err := uc.authenticationRepo.WithTx(tx.WriterDB()).ExistsByTypeAndIdentifier(ctx, authModel.AuthenticationTypeEmailPassword, m.Email)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return errors.Join(ErrRequestEmailVerificationEmailAlreadyExists, fmt.Errorf("email=[%s]", m.Email))
		}

		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to create email verification: %w", err)
		}

		lang := contexts.MustLanguage(ctx)
		minute := i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonFormatMinute, TemplateData: map[string]int{"Minute": cfg.EmailVerificationExpiresInMinute()}})
		body := i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.RegistrationEmailRequest_email_verificationBody, TemplateData: map[string]string{
			"ExpiresInMinute": minute,
			"Code":            m.Requested.PINCode,
		}})
		msg, err := message.New(ctx, job.SendEmailJob.String(), job.SendEmailPayload{
			From:    config.Email().From,
			To:      m.Email,
			Subject: i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.RegistrationEmailRequest_email_verificationTitle}),
			Body:    body,
		})
		if err != nil {
			return fmt.Errorf("failed to create worker message: %w", err)
		}
		if err := uc.producer.Do(ctx, msg); err != nil {
			return fmt.Errorf("failed to enqueue job: %w", err)
		}

		return nil
	}); err != nil {
		return RequestEmailVerificationOutput{}, err
	}
	return RequestEmailVerificationOutput{}, nil
}
