package model

import (
	"fmt"

	"github.com/mickamy/sampay/internal/lib/jwt"
)

type Session struct {
	UserID string
	Tokens jwt.Tokens
}

func NewSession(userID string) (Session, error) {
	tokens, err := jwt.New(userID)
	if err != nil {
		return Session{}, fmt.Errorf("failed to create jwt tokens: %w", err)
	}
	return Session{
		UserID: userID,
		Tokens: tokens,
	}, nil
}

func MustNewSession(userID string) Session {
	s, err := NewSession(userID)
	if err != nil {
		panic(err)
	}
	return s
}
