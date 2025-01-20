package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type GetOnboardingStepInput struct {
}

type GetOnboardingStepOutput struct {
	Step registrationModel.OnboardingStep
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type GetOnboardingStep interface {
	Do(ctx context.Context, input GetOnboardingStepInput) (GetOnboardingStepOutput, error)
}

type getOnboardingStep struct {
	reader                *database.Reader
	emailVerificationRepo authRepository.EmailVerification
	authRepo              authRepository.Authentication
	userRepo              userRepository.User
}

func NewGetOnboardingStep(
	reader *database.Reader,
	emailVerificationRepo authRepository.EmailVerification,
	authRepo authRepository.Authentication,
	userRepo userRepository.User,
) GetOnboardingStep {
	return &getOnboardingStep{
		reader:                reader,
		emailVerificationRepo: emailVerificationRepo,
		authRepo:              authRepo,
		userRepo:              userRepo,
	}
}

func (uc *getOnboardingStep) Do(ctx context.Context, input GetOnboardingStepInput) (GetOnboardingStepOutput, error) {
	var step registrationModel.OnboardingStep

	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		token := contexts.MustAnonymousUserToken(ctx)
		verification, err := uc.emailVerificationRepo.WithTx(tx.ReaderDB()).FindByVerifiedToken(ctx, token)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return errors.New("email verification not found")
		}

		auth, err := uc.authRepo.FindByTypeAndIdentifier(ctx, authModel.AuthenticationTypeEmailPassword, verification.Email)
		if err != nil {
			return fmt.Errorf("failed to find auth: %w", err)
		}
		if auth == nil {
			step = registrationModel.OnboardingStepPassword
			return nil
		}

		user, err := uc.userRepo.WithTx(tx.ReaderDB()).Get(ctx, auth.UserID, userRepository.UserJoinAttribute, userRepository.UserJoinProfile)
		if err != nil {
			return fmt.Errorf("failed to get user with attribute and profile: %w", err)
		}

		if user.Attribute.IsZero() {
			step = registrationModel.OnboardingStepAttribute
			return nil
		}

		if user.Profile.IsZero() {
			step = registrationModel.OnboardingStepProfile
			return nil
		}

		step = registrationModel.OnboardingStepCompleted

		return nil
	}); err != nil {
		return GetOnboardingStepOutput{}, err
	}

	return GetOnboardingStepOutput{Step: step}, nil
}
