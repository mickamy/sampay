package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrGetUserNotFound = commonModel.NewLocalizableError(errors.New("user not found")).
		WithMessages(i18n.Config{MessageID: i18n.UserUsecaseGet_userErrorNot_found})
)

type GetUserInput struct {
	Slug string
}

type GetUserOutput struct {
	userModel.User
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type GetUser interface {
	Do(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}

type getUser struct {
	reader   *database.Reader
	userRepo userRepository.User
}

func NewGetUser(
	reader *database.Reader,
	userRepo userRepository.User,
) GetUser {
	return &getUser{
		reader:   reader,
		userRepo: userRepo,
	}
}

func (uc *getUser) Do(ctx context.Context, input GetUserInput) (GetUserOutput, error) {
	var m *userModel.User
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		m, err = uc.userRepo.WithTx(tx.ReaderDB()).FindBySlug(ctx, input.Slug, userRepository.UserPreloadProfileAndImage, userRepository.UserPreloadLinksQRCodeAndDisplayAttributes)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if m == nil {
			return errors.Join(ErrGetUserNotFound, fmt.Errorf("slug=[%s]", input.Slug))
		}

		return nil
	}); err != nil {
		return GetUserOutput{}, err
	}

	return GetUserOutput{User: *m}, nil
}
