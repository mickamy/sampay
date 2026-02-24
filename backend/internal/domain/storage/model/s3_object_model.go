package model

import (
	"time"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type S3Object struct {
	ID        string
	Bucket    string
	Key       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
