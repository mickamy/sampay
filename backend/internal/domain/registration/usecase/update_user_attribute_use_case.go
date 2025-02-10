package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type UpdateUserAttributeInput struct {
	UsageCategoryType string
}

type UpdateUserAttributeOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserAttribute interface {
	Do(ctx context.Context, input UpdateUserAttributeInput) (UpdateUserAttributeOutput, error)
}

type updateUserAttribute struct {
	writer            *database.Writer
	userAttributeRepo userRepository.UserAttribute
}

func NewUpdateUserAttribute(
	writer *database.Writer,
	userAttributeRepo userRepository.UserAttribute,
) UpdateUserAttribute {
	return &updateUserAttribute{
		writer:            writer,
		userAttributeRepo: userAttributeRepo,
	}
}

func (uc *updateUserAttribute) Do(ctx context.Context, input UpdateUserAttributeInput) (UpdateUserAttributeOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m := userModel.UserAttribute{
			UserID:            contexts.MustAuthenticatedUserID(ctx),
			UsageCategoryType: input.UsageCategoryType,
		}
		if err := uc.userAttributeRepo.WithTx(tx.WriterDB()).Upsert(ctx, &m); err != nil {
			return fmt.Errorf("failed to persist user attribute: %w", err)
		}

		return nil
	}); err != nil {
		return UpdateUserAttributeOutput{}, err
	}

	return UpdateUserAttributeOutput{}, nil
}
