package usecase

import (
	"context"
	"fmt"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	commonRepository "mickamy.com/sampay/internal/domain/common/repository"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

type UpdateUserLinkQRCodeInput struct {
	ID     string
	QRCode *commonModel.S3Object
}

type UpdateUserLinkQRCodeOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserLinkQRCode interface {
	Do(ctx context.Context, input UpdateUserLinkQRCodeInput) (UpdateUserLinkQRCodeOutput, error)
}

type updateUserLinkQRCode struct {
	writer       *database.Writer
	userLinkRepo userRepository.UserLink
	s3Repo       commonRepository.S3Object
}

func NewUpdateUserLinkQRCode(
	writer *database.Writer,
	userLinkRepo userRepository.UserLink,
	s3Repo commonRepository.S3Object,
) UpdateUserLinkQRCode {
	return &updateUserLinkQRCode{
		writer:       writer,
		userLinkRepo: userLinkRepo,
		s3Repo:       s3Repo,
	}
}

func (uc *updateUserLinkQRCode) Do(ctx context.Context, input UpdateUserLinkQRCodeInput) (UpdateUserLinkQRCodeOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m, err := uc.userLinkRepo.WithTx(tx.WriterDB()).Find(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("failed to find user link: %w", err)
		}
		if m == nil {
			return fmt.Errorf("user link not found: id=[%s]", input.ID)
		}

		oldQRCodeID := m.QRCodeID
		m.SetQRCode(input.QRCode)

		if err := uc.userLinkRepo.WithTx(tx.WriterDB()).Update(ctx, m); err != nil {
			return fmt.Errorf("failed to update user link: %w", err)
		}

		if oldQRCodeID != nil {
			if err := uc.s3Repo.WithTx(tx.WriterDB()).Delete(ctx, *oldQRCodeID); err != nil {
				return fmt.Errorf("failed to delete s3 object: %w", err)
			}
		}
		return nil
	}); err != nil {
		return UpdateUserLinkQRCodeOutput{}, err
	}

	return UpdateUserLinkQRCodeOutput{}, nil
}
