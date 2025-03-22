package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CompleteOnboardingInput struct {
}

type CompleteOnboardingOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CompleteOnboarding interface {
	Do(ctx context.Context, input CompleteOnboardingInput) (CompleteOnboardingOutput, error)
}

type completeOnboarding struct {
	writer            *database.Writer
	userAttributeRepo userRepository.UserAttribute
}

func NewCompleteOnboarding(
	writer *database.Writer,
	userAttributeRepo userRepository.UserAttribute,
) CompleteOnboarding {
	return &completeOnboarding{
		writer:            writer,
		userAttributeRepo: userAttributeRepo,
	}
}

func (uc *completeOnboarding) Do(ctx context.Context, input CompleteOnboardingInput) (CompleteOnboardingOutput, error) {
	id := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		attr, err := uc.userAttributeRepo.WithTx(tx.WriterDB()).Find(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if attr == nil {
			return fmt.Errorf("user atttribute not found: user_id=[%s]", id)
		}

		attr.OnboardingCompleted = true
		if err := uc.userAttributeRepo.WithTx(tx.WriterDB()).Update(ctx, attr); err != nil {
			return fmt.Errorf("failed to update user attribute: %w", err)
		}

		return nil
	}); err != nil {
		return CompleteOnboardingOutput{}, err
	}

	return CompleteOnboardingOutput{}, nil
}
