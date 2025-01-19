package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/go-sqs-worker/producer"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrVerifyEmailInvalidToken = commonModel.
		NewLocalizableError(errors.New("invalid pin code")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseVerify_emailErrorInvalid_pin_code})
)

type VerifyEmailInput struct {
	Token   string
	PINCode string
}

type VerifyEmailOutput struct {
	Session authModel.Session
	Token   string
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type VerifyEmail interface {
	Do(ctx context.Context, input VerifyEmailInput) (VerifyEmailOutput, error)
}

type verifyEmail struct {
	writer                *database.Writer
	producer              *producer.Producer
	emailVerificationRepo authRepository.EmailVerification
	userRepo              userRepository.User
	sessionRepo           authRepository.Session
}

func NewVerifyEmail(
	writer *database.Writer,
	producer *producer.Producer,
	emailVerificationRepo authRepository.EmailVerification,
	userRepo userRepository.User,
	sessionRepo authRepository.Session,
) VerifyEmail {
	return &verifyEmail{
		writer:                writer,
		producer:              producer,
		emailVerificationRepo: emailVerificationRepo,
		userRepo:              userRepo,
		sessionRepo:           sessionRepo,
	}
}

func (uc *verifyEmail) Do(ctx context.Context, input VerifyEmailInput) (VerifyEmailOutput, error) {
	var session authModel.Session
	var token string
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		var err error
		verification, err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).FindByRequestedTokenAndPinCode(
			ctx,
			input.Token,
			input.PINCode,
			authRepository.EmailVerificationJoinVerified,
			authRepository.EmailVerificationNotConsumed,
		)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return ErrVerifyEmailInvalidToken
		}

		if verification.IsVerified() {
			return ErrVerifyEmailInvalidToken
		}
		if err := verification.Verify(); err != nil {
			return fmt.Errorf("failed to verify email verification: %w", err)
		}
		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Update(ctx, verification); err != nil {
			return fmt.Errorf("failed to update email verification: %w", err)
		}

		token = verification.Verified.Token

		user := userModel.User{}
		if err := uc.userRepo.WithTx(tx.WriterDB()).Create(ctx, &user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		session, err = authModel.NewSession(user.ID)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		if err := uc.sessionRepo.Create(ctx, session); err != nil {
			return fmt.Errorf("failed to persist session: %w", err)
		}

		return nil
	}); err != nil {
		return VerifyEmailOutput{}, err
	}

	return VerifyEmailOutput{Session: session, Token: token}, nil
}
