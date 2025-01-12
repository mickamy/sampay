package fixture

import (
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
)

func UserLinkProvider(setter func(m *model.UserLinkProvider)) model.UserLinkProvider {
	m := model.UserLinkProvider{
		Type:         UserLinkProviderType(),
		DisplayOrder: commonFixture.Int(),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}

func UserLinkProviderType() model.UserLinkProviderType {
	return commonFixture.RandomStringer(model.LinkProviderTypes())
}
