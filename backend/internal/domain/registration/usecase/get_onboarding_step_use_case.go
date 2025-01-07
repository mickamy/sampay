package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
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
	reader   *database.Reader
	userRepo userRepository.User
}

func NewGetOnboardingStep(
	reader *database.Reader,
	userRepo userRepository.User,
) GetOnboardingStep {
	return &getOnboardingStep{
		reader:   reader,
		userRepo: userRepo,
	}
}

func (uc *getOnboardingStep) Do(ctx context.Context, input GetOnboardingStepInput) (GetOnboardingStepOutput, error) {
	var step registrationModel.OnboardingStep

	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		user, err := uc.userRepo.WithTx(tx.ReaderDB()).Get(ctx, contexts.MustAuthenticatedUser(ctx).ID, userRepository.UserPreloadAttribute, userRepository.UserPreloadProfile)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		if user.Attribute.IsZero() {
			step = registrationModel.OnboardingStepAttribute
			return nil
		}

		if user.Profile.IsZero() {
			step = registrationModel.OnboardingStepProfile
			return nil
		}

		step = registrationModel.OnboardingStepComplete

		return nil
	}); err != nil {
		return GetOnboardingStepOutput{}, err
	}

	return GetOnboardingStepOutput{Step: step}, nil
}
