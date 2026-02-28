package fixture

import (
	"github.com/google/uuid"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func OAuthAccount(setter func(m *model.OAuthAccount)) model.OAuthAccount {
	m := model.OAuthAccount{
		ID:       ulid.New(),
		Provider: string(model.OAuthProviderLINE),
		UID:      uuid.NewString(),
	}
	if setter != nil {
		setter(&m)
	}
	return m
}
