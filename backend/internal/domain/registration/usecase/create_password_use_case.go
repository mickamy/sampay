package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CreatePasswordInput struct {
	Email    string
	Password string
}

type CreatePasswordOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreatePassword interface {
	Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error)
}

type createPassword struct {
	writer             *database.Writer
	authenticationRepo authRepository.Authentication
}

func NewCreatePassword(
	writer *database.Writer,
	authenticationRepo authRepository.Authentication,
) CreatePassword {
	return &createPassword{
		writer:             writer,
		authenticationRepo: authenticationRepo,
	}
}

func (uc *createPassword) Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m, err := authModel.NewAuthenticationEmailPassword(contexts.MustAuthenticatedUserID(ctx), input.Email, input.Password)
		if err != nil {
			return fmt.Errorf("failed to create authentication model: %w", err)
		}

		if err := uc.authenticationRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to persist user attribute: %w", err)
		}

		return nil
	}); err != nil {
		return CreatePasswordOutput{}, err
	}

	return CreatePasswordOutput{}, nil
}
