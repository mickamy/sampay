package model

import (
	"time"

	amodel "github.com/mickamy/sampay/internal/domain/auth/model"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type EndUser struct {
	UserID    string `db:",primaryKey" map:"Id"`
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time

	OAuthAccounts []amodel.OAuthAccount `rel:"has_many,foreign_key:end_user_id"`
}
