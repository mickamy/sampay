package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type GetMeInput struct {
}

type GetMeOutput struct {
	userModel.User
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type GetMe interface {
	Do(ctx context.Context, input GetMeInput) (GetMeOutput, error)
}

type getMe struct {
	reader   *database.Reader
	userRepo userRepository.User
}

func NewGetMe(
	reader *database.Reader,
	userRepo userRepository.User,
) GetMe {
	return &getMe{
		reader:   reader,
		userRepo: userRepo,
	}
}

func (uc *getMe) Do(ctx context.Context, input GetMeInput) (GetMeOutput, error) {
	var m userModel.User
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		m, err = uc.userRepo.WithTx(tx.ReaderDB()).Get(ctx, contexts.MustAuthenticatedUserID(ctx), userRepository.UserPreloadProfileAndImage, userRepository.UserPreloadLinksQRCodeAndDisplayAttributes)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		return nil
	}); err != nil {
		return GetMeOutput{}, err
	}

	return GetMeOutput{User: m}, nil
}
