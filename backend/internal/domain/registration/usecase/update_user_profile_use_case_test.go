package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/ptr"
)

func TestCreateUserProfile_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).UpdateUserProfile
	_, err := sut.Do(ctx, usecase.UpdateUserProfileInput{
		Name:  gofakeit.GlobalFaker.Name(),
		Slug:  gofakeit.GlobalFaker.Username(),
		Bio:   ptr.Of(gofakeit.GlobalFaker.Sentence(20)),
		Image: ptr.Of(commonFixture.S3Object(nil)),
	})

	// assert
	require.NoError(t, err)
}
