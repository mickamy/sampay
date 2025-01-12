package fixture

import (
	"mickamy.com/sampay/internal/domain/user/model"
)

func UserLink(setter func(m *model.UserLink)) model.UserLink {
	m := model.UserLink{
		ProviderType: UserLinkProviderType(),
	}

	switch m.ProviderType {
	case model.UserLinkProviderTypeKyash:
		m.URI = UserLinkURIKyash()
	case model.UserLinkProviderTypePayPay:
		m.URI = UserLinkURIPayPay()
	case model.UserLinkProviderTypeAmazon:
		m.URI = UserLinkURIAmazon()
	default:
		m.URI = UserLinkURIOther()
	}

	if setter != nil {
		setter(&m)
	}

	return m
}

func UserLinkURIKyash() string {
	return "kyash://qr/u/1234567890"
}

func UserLinkURIPayPay() string {
	return "https://qr.paypay.ne.jp/1234567890"
}

func UserLinkURIAmazon() string {
	return "https://www.amazon.co.jp/hz/wishlist/ls/1234567890"
}

func UserLinkURIOther() string {
	return "https://example.com"
}
