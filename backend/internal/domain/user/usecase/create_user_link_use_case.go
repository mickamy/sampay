package usecase

import (
	"context"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
)

type CreateUserLinkInput struct {
	userModel.UserLink
}

type CreateUserLinkOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateUserLink interface {
	Do(ctx context.Context, input CreateUserLinkInput) (CreateUserLinkOutput, error)
}

type createUserLink struct {
	writer       *database.Writer
	userLinkRepo userRepository.UserLink
}

func NewCreateUserLink(
	writer *database.Writer,
	userLinkRepo userRepository.UserLink,
) CreateUserLink {
	return &createUserLink{
		writer:       writer,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *createUserLink) Do(ctx context.Context, input CreateUserLinkInput) (CreateUserLinkOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Create(ctx, &input.UserLink); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return CreateUserLinkOutput{}, err
	}

	return CreateUserLinkOutput{}, nil
}
