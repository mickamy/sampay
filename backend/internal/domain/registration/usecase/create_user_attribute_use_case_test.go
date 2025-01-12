package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestCreateUserAttribute_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreateUserAttribute
	_, err := sut.Do(ctx, usecase.CreateUserAttributeInput{
		UsageCategoryType: "other",
	})

	// assert
	require.NoError(t, err)
}
