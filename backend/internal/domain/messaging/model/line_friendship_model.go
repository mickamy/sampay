package model

import (
	"time"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type LineFriendship struct {
	EndUserID string `db:",primaryKey"`
	IsFriend  bool
	UpdatedAt time.Time
}
