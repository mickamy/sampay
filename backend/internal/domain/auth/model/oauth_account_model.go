package model

import (
	"time"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type OAuthAccount struct {
	ID        string
	EndUserID string
	Provider  string
	UID       string
	CreatedAt time.Time
}
