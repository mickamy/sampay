package repository

import (
	"mickamy.com/sampay/internal/cli/infra/storage/database"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserLink interface {
	WithTx(tx *database.DB) UserLink
}

type userLink struct {
	db *database.DB
}

func NewUserLink(db *database.DB) UserLink {
	return &userLink{db: db}
}

func (repo *userLink) WithTx(tx *database.DB) UserLink {
	return &userLink{db: tx}
}
