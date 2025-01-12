package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
)

var (
	ErrUpdateUserLinkNotFound = errors.New("user link not found")
)

type UpdateUserLinkInput struct {
	ID           string
	ProviderType *userModel.UserLinkProviderType
	URI          *string
	Name         *string
	DisplayOrder *int
}

type UpdateUserLinkOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserLink interface {
	Do(ctx context.Context, input UpdateUserLinkInput) (UpdateUserLinkOutput, error)
}

type updateUserLink struct {
	writer       *database.Writer
	userLinkRepo userRepository.UserLink
}

func NewUpdateUserLink(
	writer *database.Writer,
	userLinkRepo userRepository.UserLink,
) UpdateUserLink {
	return &updateUserLink{
		writer:       writer,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *updateUserLink) Do(ctx context.Context, input UpdateUserLinkInput) (UpdateUserLinkOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m, err := uc.userLinkRepo.WithTx(tx.WriterDB()).Find(ctx, input.ID, userRepository.UserLinkJoinDisplayAttribute)
		if err != nil {
			return fmt.Errorf("failed to find user link: %w", err)
		}
		if m == nil {
			return errors.Join(ErrUpdateUserLinkNotFound, fmt.Errorf("id=[%s]", input.ID))
		}

		if input.ProviderType != nil {
			m.ProviderType = *input.ProviderType
		}
		if input.URI != nil {
			m.URI = *input.URI
		}
		if input.Name != nil {
			m.DisplayAttribute.Name = *input.Name
		}
		if input.DisplayOrder != nil {
			m.DisplayAttribute.DisplayOrder = *input.DisplayOrder
		}

		if err := uc.userLinkRepo.WithTx(tx.WriterDB().FullSaveAssociations()).Update(ctx, m); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return UpdateUserLinkOutput{}, err
	}

	return UpdateUserLinkOutput{}, nil
}
