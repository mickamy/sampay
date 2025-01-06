package dto

import (
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	"mickamy.com/sampay/internal/lib/jwt"
)

func NewTokens(tokens jwt.Tokens) *authv1.Tokens {
	return &authv1.Tokens{
		Access:  newToken(tokens.Access),
		Refresh: newToken(tokens.Refresh),
	}
}

func newToken(token jwt.Token) *authv1.Token {
	return &authv1.Token{
		Value:     token.Value,
		ExpiresAt: timestamppb.New(token.ExpiresAt),
	}
}
