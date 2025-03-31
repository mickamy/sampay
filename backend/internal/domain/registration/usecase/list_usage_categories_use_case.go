package usecase

import (
	"context"
	"fmt"

	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	registrationRepository "mickamy.com/sampay/internal/domain/registration/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

type ListUsageCategoriesInput struct {
}

type ListUsageCategoriesOutput struct {
	Categories []registrationModel.UsageCategory
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type ListUsageCategories interface {
	Do(ctx context.Context, input ListUsageCategoriesInput) (ListUsageCategoriesOutput, error)
}

type listUsageCategories struct {
	reader            *database.Reader
	usageCategoryRepo registrationRepository.UsageCategory
}

func NewListUsageCategories(
	reader *database.Reader,
	usageCategoryRepo registrationRepository.UsageCategory,
) ListUsageCategories {
	return &listUsageCategories{
		reader:            reader,
		usageCategoryRepo: usageCategoryRepo,
	}
}

func (uc *listUsageCategories) Do(ctx context.Context, input ListUsageCategoriesInput) (ListUsageCategoriesOutput, error) {
	var categories []registrationModel.UsageCategory

	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		categories, err = uc.usageCategoryRepo.WithTx(tx.ReaderDB()).List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list usage categories: %w", err)
		}

		return nil
	}); err != nil {
		return ListUsageCategoriesOutput{}, err
	}

	return ListUsageCategoriesOutput{Categories: categories}, nil
}
