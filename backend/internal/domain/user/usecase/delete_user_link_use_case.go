package usecase

import (
	"context"

	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

type DeleteUserLinkInput struct {
	ID string
}

type DeleteUserLinkOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type DeleteUserLink interface {
	Do(ctx context.Context, input DeleteUserLinkInput) (DeleteUserLinkOutput, error)
}

type deleteUserLink struct {
	writer       *database.Writer
	userLinkRepo userRepository.UserLink
}

func NewDeleteUserLink(
	writer *database.Writer,
	userLinkRepo userRepository.UserLink,
) DeleteUserLink {
	return &deleteUserLink{
		writer:       writer,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *deleteUserLink) Do(ctx context.Context, input DeleteUserLinkInput) (DeleteUserLinkOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Delete(ctx, input.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return DeleteUserLinkOutput{}, err
	}

	return DeleteUserLinkOutput{}, nil
}
