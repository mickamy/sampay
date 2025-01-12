package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CreateUserAttributeInput struct {
	UsageCategoryType string
}

type CreateUserAttributeOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateUserAttribute interface {
	Do(ctx context.Context, input CreateUserAttributeInput) (CreateUserAttributeOutput, error)
}

type createUserAttribute struct {
	writer            *database.Writer
	userAttributeRepo userRepository.UserAttribute
}

func NewCreateUserAttribute(
	writer *database.Writer,
	userAttributeRepo userRepository.UserAttribute,
) CreateUserAttribute {
	return &createUserAttribute{
		writer:            writer,
		userAttributeRepo: userAttributeRepo,
	}
}

func (uc *createUserAttribute) Do(ctx context.Context, input CreateUserAttributeInput) (CreateUserAttributeOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m := userModel.UserAttribute{
			UserID:            contexts.MustAuthenticatedUserID(ctx),
			UsageCategoryType: input.UsageCategoryType,
		}
		if err := uc.userAttributeRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to persist user attribute: %w", err)
		}

		return nil
	}); err != nil {
		return CreateUserAttributeOutput{}, err
	}

	return CreateUserAttributeOutput{}, nil
}
