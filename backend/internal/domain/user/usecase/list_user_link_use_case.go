package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
)

type ListUserLinkInput struct {
	UserID string
}

type ListUserLinkOutput struct {
	Links []userModel.UserLink
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type ListUserLink interface {
	Do(ctx context.Context, input ListUserLinkInput) (ListUserLinkOutput, error)
}

type listUserLink struct {
	reader       *database.Reader
	userLinkRepo userRepository.UserLink
}

func NewListUserLink(
	reader *database.Reader,
	userLinkRepo userRepository.UserLink,
) ListUserLink {
	return &listUserLink{
		reader:       reader,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *listUserLink) Do(ctx context.Context, input ListUserLinkInput) (ListUserLinkOutput, error) {
	var ms []userModel.UserLink
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		ms, err = uc.userLinkRepo.WithTx(tx.ReaderDB()).ListByUserID(ctx, input.UserID, userRepository.UserLinkJoinDisplayAttribute)
		if err != nil {
			return fmt.Errorf("failed to list user links: %w", err)
		}

		return nil
	}); err != nil {
		return ListUserLinkOutput{}, err
	}

	return ListUserLinkOutput{Links: ms}, nil
}
