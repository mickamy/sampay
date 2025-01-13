package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	commonRepository "mickamy.com/sampay/internal/domain/common/repository"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type DeleteUserProfileImageInput struct {
}

type DeleteUserProfileImageOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type DeleteUserProfileImage interface {
	Do(ctx context.Context, input DeleteUserProfileImageInput) (DeleteUserProfileImageOutput, error)
}

type deleteUserProfileImage struct {
	writer          *database.Writer
	userProfileRepo userRepository.UserProfile
	s3ObjectRepo    commonRepository.S3Object
}

func NewDeleteUserProfileImage(
	writer *database.Writer,
	userProfileRepo userRepository.UserProfile,
	s3ObjectRepo commonRepository.S3Object,
) DeleteUserProfileImage {
	return &deleteUserProfileImage{
		writer:          writer,
		userProfileRepo: userProfileRepo,
		s3ObjectRepo:    s3ObjectRepo,
	}
}

func (uc *deleteUserProfileImage) Do(ctx context.Context, input DeleteUserProfileImageInput) (DeleteUserProfileImageOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		userID := contexts.MustAuthenticatedUserID(ctx)
		profile, err := uc.userProfileRepo.WithTx(tx.WriterDB()).Find(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to find user profile: %w", err)
		}
		if profile == nil {
			return fmt.Errorf("user profile not found: id=[%s]", userID)
		}

		if profile.ImageID == nil {
			return nil
		}

		imageID := *profile.ImageID
		profile.SetImage(nil)
		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Update(ctx, profile); err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}

		if err := uc.s3ObjectRepo.WithTx(tx.WriterDB()).Delete(ctx, imageID); err != nil {
			return fmt.Errorf("failed to delete user profile image: %w", err)
		}

		return nil
	}); err != nil {
		return DeleteUserProfileImageOutput{}, err
	}

	return DeleteUserProfileImageOutput{}, nil
}
