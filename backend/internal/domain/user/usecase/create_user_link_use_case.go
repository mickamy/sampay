package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
)

type CreateUserLinkInput struct {
	ProviderType userModel.UserLinkProviderType
	URI          string
	Name         string
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
	m := userModel.UserLink{
		ProviderType: input.ProviderType,
		URI:          input.URI,
		DisplayAttribute: userModel.UserLinkDisplayAttribute{
			Name: input.Name,
		},
	}
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		order, err := uc.userLinkRepo.WithTx(tx.WriterDB()).GetLastDisplayOrderByUserID(ctx, m.UserID)
		if err != nil {
			return fmt.Errorf("failed to get last display order by user id: %w", err)
		}

		m.DisplayAttribute.DisplayOrder = order + 1

		if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return CreateUserLinkOutput{}, err
	}

	return CreateUserLinkOutput{}, nil
}
