package model

import (
	"time"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type User struct {
	ID        string
	CreatedAt time.Time
}
