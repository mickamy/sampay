package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type UserLink struct {
	ID           string
	ProviderType userModel.UserLinkProviderType
	Name         string
	URI          string
	QRCode       *commonModel.S3Object
}

func (l UserLink) AsModel() userModel.UserLink {
	return userModel.UserLink{
		ID:           l.ID,
		ProviderType: l.ProviderType,
		URI:          l.URI,
		QRCode:       l.QRCode,
		DisplayAttribute: userModel.UserLinkDisplayAttribute{
			UserLinkID: l.ID,
			Name:       l.Name,
		},
	}
}

type UpdateUserLinksInput struct {
	UserLinks []UserLink
}

type UpdateUserLinksOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserLinks interface {
	Do(ctx context.Context, input UpdateUserLinksInput) (UpdateUserLinksOutput, error)
}

type updateUserLinks struct {
	writer       *database.Writer
	userLinkRepo userRepository.UserLink
}

func NewUpdateUserLinks(
	writer *database.Writer,
	userLinkRepo userRepository.UserLink,
) UpdateUserLinks {
	return &updateUserLinks{
		writer:       writer,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *updateUserLinks) Do(ctx context.Context, input UpdateUserLinksInput) (UpdateUserLinksOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		existing, err := uc.userLinkRepo.WithTx(tx.WriterDB()).ListByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to find user link: %w", err)
		}

		var deleted []userModel.UserLink
		for _, e := range existing {
			for _, i := range input.UserLinks {
				if e.ID == i.ID {
					deleted = append(deleted, e)
					break
				}
			}
		}

		for _, d := range deleted {
			if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Delete(ctx, d.ID); err != nil {
				return fmt.Errorf("failed to delete user link: %w", err)
			}
		}

		for _, i := range input.UserLinks {
			m := i.AsModel()
			m.UserID = userID
			if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Upsert(ctx, &m); err != nil {
				return fmt.Errorf("failed to update user link: %w", err)
			}
		}

		return nil
	}); err != nil {
		return UpdateUserLinksOutput{}, err
	}

	return UpdateUserLinksOutput{}, nil
}
