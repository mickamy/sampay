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

type UpdateUserProfileInput struct {
	Name  string
	Slug  string
	Bio   *string
	Image *commonModel.S3Object
}

type UpdateUserProfileOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserProfile interface {
	Do(ctx context.Context, input UpdateUserProfileInput) (UpdateUserProfileOutput, error)
}

type updateUserProfile struct {
	writer          *database.Writer
	userRepo        userRepository.User
	userProfileRepo userRepository.UserProfile
}

func NewUpdateUserProfile(
	writer *database.Writer,
	userRepo userRepository.User,
	userProfileRepo userRepository.UserProfile,
) UpdateUserProfile {
	return &updateUserProfile{
		writer:          writer,
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
	}
}

func (uc *updateUserProfile) Do(ctx context.Context, input UpdateUserProfileInput) (UpdateUserProfileOutput, error) {
	id := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		user, err := uc.userRepo.WithTx(tx.WriterDB()).Find(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if user == nil {
			return fmt.Errorf("user not found: id=[%s]", id)
		}

		user.Slug = input.Slug
		if err := uc.userRepo.WithTx(tx.WriterDB()).Update(ctx, user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		profile := userModel.UserProfile{
			UserID: id,
			Name:   input.Name,
			Bio:    input.Bio,
		}
		profile.SetImage(input.Image)

		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Upsert(ctx, &profile); err != nil {
			return fmt.Errorf("failed to persist user profile: %w", err)
		}

		return nil
	}); err != nil {
		return UpdateUserProfileOutput{}, err
	}

	return UpdateUserProfileOutput{}, nil
}
