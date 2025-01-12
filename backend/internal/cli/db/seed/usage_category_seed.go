package seed

import (
	"context"
	"fmt"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	registrationRepository "mickamy.com/sampay/internal/domain/registration/repository"
)

func seedUsageCategory(ctx context.Context, writer *database.Writer, env config.Env) error {
	categories := []registrationModel.UsageCategory{
		{Type: "business", DisplayOrder: 1},
		{Type: "influencer", DisplayOrder: 2},
		{Type: "personal", DisplayOrder: 3},
		{Type: "entertainment", DisplayOrder: 4},
		{Type: "fashion", DisplayOrder: 5},
		{Type: "restaurant", DisplayOrder: 6},
		{Type: "health", DisplayOrder: 7},
		{Type: "non_profit", DisplayOrder: 8},
		{Type: "tech", DisplayOrder: 9},
		{Type: "tourism", DisplayOrder: 10},
		{Type: "other", DisplayOrder: 99},
	}

	if err := writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		repo := registrationRepository.NewUsageCategory(tx.WriterDB())
		for _, category := range categories {
			if err := repo.Upsert(ctx, &category); err != nil {
				return fmt.Errorf("failed to upsert usage category [%+v]: %w", category, err)
			}
		}
		return nil
	}); err != nil {
		return err

	}

	return nil
}
