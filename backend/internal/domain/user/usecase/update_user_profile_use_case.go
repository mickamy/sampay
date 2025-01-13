package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type UpdateUserProfileInput struct {
	Name string
	Bio  *string
}

type UpdateUserProfileOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UpdateUserProfile interface {
	Do(ctx context.Context, input UpdateUserProfileInput) (UpdateUserProfileOutput, error)
}

type updateUserProfile struct {
	writer          *database.Writer
	userProfileRepo userRepository.UserProfile
}

func NewUpdateUserProfile(
	writer *database.Writer,
	userProfileRepo userRepository.UserProfile,
) UpdateUserProfile {
	return &updateUserProfile{
		writer:          writer,
		userProfileRepo: userProfileRepo,
	}
}

func (uc *updateUserProfile) Do(ctx context.Context, input UpdateUserProfileInput) (UpdateUserProfileOutput, error) {
	m := userModel.UserProfile{
		UserID: contexts.MustAuthenticatedUserID(ctx),
		Name:   input.Name,
	}
	if input.Bio != nil {
		m.Bio = input.Bio
	}
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Update(ctx, &m); err != nil {
			return fmt.Errorf("failed to update user profile: %w", err)
		}
		return nil
	}); err != nil {
		return UpdateUserProfileOutput{}, err
	}

	return UpdateUserProfileOutput{}, nil
}
