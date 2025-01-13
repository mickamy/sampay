package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	commonRepository "mickamy.com/sampay/internal/domain/common/repository"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type UpdateUserProfileImageInput struct {
	Image *commonModel.S3Object
}

type UpdateUserProfileImageOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserProfileImage interface {
	Do(ctx context.Context, input UpdateUserProfileImageInput) (UpdateUserProfileImageOutput, error)
}

type updateUserProfileImage struct {
	writer          *database.Writer
	userProfileRepo userRepository.UserProfile
	s3Repo          commonRepository.S3Object
}

func NewUpdateUserProfileImage(
	writer *database.Writer,
	userProfileRepo userRepository.UserProfile,
	s3Repo commonRepository.S3Object,
) UpdateUserProfileImage {
	return &updateUserProfileImage{
		writer:          writer,
		userProfileRepo: userProfileRepo,
		s3Repo:          s3Repo,
	}
}

func (uc *updateUserProfileImage) Do(ctx context.Context, input UpdateUserProfileImageInput) (UpdateUserProfileImageOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		userID := contexts.MustAuthenticatedUserID(ctx)
		m, err := uc.userProfileRepo.WithTx(tx.WriterDB()).Find(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to find user profile: %w", err)
		}
		if m == nil {
			return fmt.Errorf("user profile not found: user_id=[%s]", userID)
		}

		oldImageID := m.ImageID

		m.SetImage(input.Image)
		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Update(ctx, m); err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}

		if oldImageID != nil {
			if err := uc.s3Repo.WithTx(tx.WriterDB()).Delete(ctx, *oldImageID); err != nil {
				return fmt.Errorf("failed to delete s3 object: %w", err)
			}
		}
		return nil
	}); err != nil {
		return UpdateUserProfileImageOutput{}, err
	}

	return UpdateUserProfileImageOutput{}, nil
}
