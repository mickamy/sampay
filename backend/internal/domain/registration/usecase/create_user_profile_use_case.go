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

type CreateUserProfileInput struct {
	Name  string
	Slug  string
	Bio   *string
	Image *commonModel.S3Object
}

type CreateUserProfileOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateUserProfile interface {
	Do(ctx context.Context, input CreateUserProfileInput) (CreateUserProfileOutput, error)
}

type createUserProfile struct {
	writer          *database.Writer
	userRepo        userRepository.User
	userProfileRepo userRepository.UserProfile
}

func NewCreateUserProfile(
	writer *database.Writer,
	userRepo userRepository.User,
	userProfileRepo userRepository.UserProfile,
) CreateUserProfile {
	return &createUserProfile{
		writer:          writer,
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
	}
}

func (uc *createUserProfile) Do(ctx context.Context, input CreateUserProfileInput) (CreateUserProfileOutput, error) {
	id := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		user := userModel.User{
			ID:   id,
			Slug: input.Slug,
		}
		if err := uc.userRepo.WithTx(tx.WriterDB()).Update(ctx, &user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		profile := userModel.UserProfile{
			UserID: id,
			Name:   input.Name,
			Bio:    input.Bio,
		}
		profile.SetImage(input.Image)

		if err := uc.userProfileRepo.WithTx(tx.WriterDB()).Create(ctx, &profile); err != nil {
			return fmt.Errorf("failed to persist user profile: %w", err)
		}

		return nil
	}); err != nil {
		return CreateUserProfileOutput{}, err
	}

	return CreateUserProfileOutput{}, nil
}
