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
		userID, err := contexts.AuthenticatedUserID(ctx)
		if userID != "" {
			// check if user has password
			auths, err := uc.authRepo.WithTx(tx.ReaderDB()).FindByUserIDAndType(ctx, userID, authModel.AuthenticationTypePassword)
			if err != nil {
				return fmt.Errorf("failed to find auth: %w", err)
			}
			if auths == nil {
				step = registrationModel.OnboardingStepPassword
				return nil
			}
		} else if err != nil {
			// authenticated by anonymous token
			token, err := contexts.AnonymousUserToken(ctx)
			if err != nil {
				return fmt.Errorf("failed to get user id: %w", err)
			}

			// get email from email verification
			verification, err := uc.emailVerificationRepo.WithTx(tx.ReaderDB()).FindByVerifiedToken(ctx, token)
			if err != nil {
				return fmt.Errorf("failed to find email verification: %w", err)
			}
			if verification == nil {
				return errors.New("email verification not found")
			}

			// check if user has password
			auth, err := uc.authRepo.FindByTypeAndIdentifier(ctx, authModel.AuthenticationTypePassword, verification.Email)
			if err != nil {
				return fmt.Errorf("failed to find auth: %w", err)
			}
			if auth == nil {
				user, err := uc.userRepo.FindByEmail(ctx, verification.Email)
				if err != nil {
					return fmt.Errorf("failed to find user by email: %w", err)
				}
				if user == nil {
					return fmt.Errorf("user not found: email=[%s]", verification.Email)
				}
				userID = user.ID
			} else {
				userID = auth.UserID
			}
		}

		user, err := uc.userRepo.WithTx(tx.ReaderDB()).Get(ctx, userID, userRepository.UserJoinAttribute, userRepository.UserJoinProfile)
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
