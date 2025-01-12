package seed

import (
	"context"
	"fmt"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
)

func seedUserLinkProvider(ctx context.Context, writer *database.Writer, env config.Env) error {
	providers := []userModel.UserLinkProvider{
		{Type: userModel.UserLinkProviderTypeKyash, DisplayOrder: 1},
		{Type: userModel.UserLinkProviderTypePayPay, DisplayOrder: 2},
		{Type: userModel.UserLinkProviderTypeAmazon, DisplayOrder: 3},
		{Type: userModel.UserLinkProviderTypeOther, DisplayOrder: 99},
	}

	if err := writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		repo := userRepository.NewUserLinkProvider(tx.WriterDB())
		for _, provider := range providers {
			if err := repo.Upsert(ctx, &provider); err != nil {
				return fmt.Errorf("failed to upsert user link provider [%+v]: %w", provider, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
