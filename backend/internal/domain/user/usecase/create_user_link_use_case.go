package usecase

import (
	"context"
	"fmt"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
	"mickamy.com/sampay/internal/lib/aws/s3"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CreateUserLinkInput struct {
	ProviderType userModel.UserLinkProviderType
	URI          string
	Name         string
	QRCode       *commonModel.S3Object
}

type CreateUserLinkOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateUserLink interface {
	Do(ctx context.Context, input CreateUserLinkInput) (CreateUserLinkOutput, error)
}

type createUserLink struct {
	writer       *database.Writer
	s3           s3.Client
	userLinkRepo userRepository.UserLink
}

func NewCreateUserLink(
	writer *database.Writer,
	s3 s3.Client,
	userLinkRepo userRepository.UserLink,
) CreateUserLink {
	return &createUserLink{
		writer:       writer,
		s3:           s3,
		userLinkRepo: userLinkRepo,
	}
}

func (uc *createUserLink) Do(ctx context.Context, input CreateUserLinkInput) (CreateUserLinkOutput, error) {
	m := userModel.UserLink{
		UserID:       contexts.MustAuthenticatedUserID(ctx),
		ProviderType: input.ProviderType,
		URI:          input.URI,
		DisplayAttribute: userModel.UserLinkDisplayAttribute{
			Name: input.Name,
		},
	}
	m.SetQRCode(input.QRCode)
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
