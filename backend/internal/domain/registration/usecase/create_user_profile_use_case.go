package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CreateUserProfileInput struct {
	Name string
	Bio  *string
}

type CreateUserProfileOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateUserProfile interface {
	Do(ctx context.Context, input CreateUserProfileInput) (CreateUserProfileOutput, error)
}

type createUserProfile struct {
	writer          *database.Writer
	userProfileRepo userRepository.UserProfile
}

func NewCreateUserProfile(
	writer *database.Writer,
	userProfileRepo userRepository.UserProfile,
) CreateUserProfile {
	return &createUserProfile{
		writer:          writer,
		userProfileRepo: userProfileRepo,
	}
}

func (uc *createUserProfile) Do(ctx context.Context, input CreateUserProfileInput) (CreateUserProfileOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		m := userModel.UserProfile{
			UserID: contexts.MustAuthenticatedUserID(ctx),
			Name:   input.Name,
			Bio:    input.Bio,
		}
		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to persist user profile: %w", err)
		}

		return nil
	}); err != nil {
		return CreateUserProfileOutput{}, err
	}

	return CreateUserProfileOutput{}, nil
}
