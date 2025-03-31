package paging

import (
	"mickamy.com/sampay/internal/infra/storage/database"
)

type Page struct {
	Index int
	Limit int
}

func New(index, limit int) Page {
	return Page{
		Index: index,
		Limit: limit,
	}
}

func (p Page) Next(total int) NextPage {
	return NextPage{
		Index: p.Index + 1,
		Limit: p.Limit,
		Total: total,
	}
}

func (p Page) Scope() database.Scope {
	return func(db *database.DB) *database.DB {
		return &database.DB{DB: db.Offset(p.Index * p.Limit).Limit(p.Limit)}
	}
}

type NextPage struct {
	Index int
	Limit int
	Total int
}
