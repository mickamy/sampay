package model

import (
	"fmt"
	"regexp"
)

type UserLinkProviderType string

const (
	UserLinkProviderTypeKyash  UserLinkProviderType = "kyash"
	UserLinkProviderTypePayPay UserLinkProviderType = "paypay"
	UserLinkProviderTypeAmazon UserLinkProviderType = "amazon"
	UserLinkProviderTypeOther  UserLinkProviderType = "other"
)

var (
	UserLinkProviderTypeKyashRegexp  = regexp.MustCompile(`^kyash://qr/u/[a-zA-Z0-9]+$`)
	UserLinkProviderTypePayPayRegexp = regexp.MustCompile(`^https://qr\.paypay\.ne\.jp/[a-zA-Z0-9_]+$`)
	UserLinkProviderTypeAmazonRegexp = regexp.MustCompile(`^https://www\.amazon(\.co)?\.jp/(hz/)?wishlist/ls/.+$`)
	UserLinkProviderTypeOtherRegexp  = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://[^/]+\..+`)
)

func (m UserLinkProviderType) String() string {
	return string(m)
}
func (m UserLinkProviderType) Regexp() *regexp.Regexp {
	switch m {
	case UserLinkProviderTypeKyash:
		return UserLinkProviderTypeKyashRegexp
	case UserLinkProviderTypePayPay:
		return UserLinkProviderTypePayPayRegexp
	case UserLinkProviderTypeAmazon:
		return UserLinkProviderTypeAmazonRegexp
	default:
		return UserLinkProviderTypeOtherRegexp
	}
}

func (m UserLinkProviderType) MatchString(s string) bool {
	return m.Regexp().MatchString(s)
}

func LinkProviderTypes() []UserLinkProviderType {
	return []UserLinkProviderType{
		UserLinkProviderTypeKyash,
		UserLinkProviderTypePayPay,
		UserLinkProviderTypeAmazon,
		UserLinkProviderTypeOther,
	}
}

func ParseLinkProviderTypeURI(uri string) UserLinkProviderType {
	for _, t := range LinkProviderTypes() {
		if t.Regexp().MatchString(uri) {
			return t
		}
	}
	return UserLinkProviderTypeOther
}

func NewLinkProviderType(s string) (UserLinkProviderType, error) {
	for _, t := range LinkProviderTypes() {
		if t.String() == s {
			return t, nil
		}
	}
	return "", fmt.Errorf("invalid link provider type: %s", s)
}

func MustNewLinkProviderType(s string) UserLinkProviderType {
	t, err := NewLinkProviderType(s)
	if err != nil {
		panic(err)
	}
	return t
}
