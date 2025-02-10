package fixture

import (
	"mickamy.com/sampay/internal/domain/user/model"
)

func UserAttribute(setter func(m *model.UserAttribute)) model.UserAttribute {
	m := model.UserAttribute{
		UsageCategoryType:   "other",
		OnboardingCompleted: true,
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
