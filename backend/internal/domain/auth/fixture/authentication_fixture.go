package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/passwd"
)

func Authentication(setter func(m *model.Authentication)) model.Authentication {
	m := model.Authentication{}

	if setter != nil {
		setter(&m)
	}

	return m
}

func AuthenticationEmailPassword(setter func(m *model.Authentication)) model.Authentication {
	m := Authentication(func(m *model.Authentication) {
		m.Type = model.AuthenticationTypeEmailPassword
		m.Identifier = gofakeit.GlobalFaker.Email()
		m.Secret = either.Must(passwd.New("P@ssw0rd", 16))
	})

	if setter != nil {
		setter(&m)
	}

	return m
}
