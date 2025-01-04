package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mickamy.com/sampay/config"
)

const (
	day  = time.Hour * 24
	week = day * 7

	accessTokenExpiresIn  = time.Hour
	refreshTokenExpiresIn = week * 2

	idKey = "id"

	expiredTokenErrorMessage = "Token is expired"
)

var (
	ErrExpiredToken = errors.New("token is expired")

	signingMethod = jwt.SigningMethodHS256
	signingSecret = config.Auth().SigningSecretBytes()
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

func (t Token) Expiration() time.Duration {
	return time.Until(t.ExpiresAt)
}

func (t Token) IsEmpty() bool {
	return t.Value == ""
}

type Tokens struct {
	Access  Token
	Refresh Token
}

func (ts Tokens) IsEmpty() bool {
	return ts.Access.IsEmpty() && ts.Refresh.IsEmpty()
}

func New(id string) (Tokens, error) {
	access, err := generateAccessToken(id)
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := generateRefreshToken(id, access)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{Access: access, Refresh: refresh}, nil
}

func generateAccessToken(id string) (Token, error) {
	claims := jwt.MapClaims{}
	claims[idKey] = id
	exp := time.Now().Add(accessTokenExpiresIn)
	claims["exp"] = exp.Unix()

	jwtToken := jwt.NewWithClaims(signingMethod, claims)
	accessTokenValue, err := jwtToken.SignedString(signingSecret)
	if err != nil {
		return Token{}, fmt.Errorf("failed to sign access token jwt: %w", err)
	}
	return Token{Value: accessTokenValue, ExpiresAt: exp}, nil
}

func generateRefreshToken(id string, accessToken Token) (Token, error) {
	claims := jwt.MapClaims{}
	claims[idKey] = id
	claims["jwt"] = accessToken.Value
	exp := time.Now().Add(refreshTokenExpiresIn)
	claims["exp"] = exp.Unix()

	jwtToken := jwt.NewWithClaims(signingMethod, claims)
	refreshTokenValue, err := jwtToken.SignedString(signingSecret)
	if err != nil {
		return Token{}, fmt.Errorf("failed to sign refresh token jwt: %w", err)
	}
	return Token{Value: refreshTokenValue, ExpiresAt: exp}, nil
}

func Verify(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingSecret, nil
	})

	if err != nil {
		if err.Error() == expiredTokenErrorMessage {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("failed to parse jwt: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid jwt")
}

func ExtractID(tokenStr string) (string, error) {
	claims, err := Verify(tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to verify jwt: %w", err)
	}
	id, ok := claims[idKey].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract id from jwt")
	}
	return id, nil
}

func IsRefreshTokenClaims(claims jwt.MapClaims) bool {
	if _, ok := claims["jwt"].(string); ok {
		return true
	}
	return false
}
