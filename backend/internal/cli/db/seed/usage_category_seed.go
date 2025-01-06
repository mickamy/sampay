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
		{CategoryType: "business", DisplayOrder: 1},
		{CategoryType: "influencer", DisplayOrder: 2},
		{CategoryType: "personal", DisplayOrder: 3},
		{CategoryType: "entertainment", DisplayOrder: 4},
		{CategoryType: "fashion", DisplayOrder: 5},
		{CategoryType: "restaurant", DisplayOrder: 6},
		{CategoryType: "health", DisplayOrder: 7},
		{CategoryType: "non_profit", DisplayOrder: 8},
		{CategoryType: "tech", DisplayOrder: 9},
		{CategoryType: "tourism", DisplayOrder: 10},
		{CategoryType: "other", DisplayOrder: 99},
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
		return fmt.Errorf("failed to seed usage categories: %w", err)
	}

	return nil
}
