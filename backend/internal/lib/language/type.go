package language

import (
	"golang.org/x/text/language"
)

type Type string

const (
	Japanese = Type("ja")
)

func (l Type) Tag() language.Tag {
	switch l {
	case Japanese:
		return language.Japanese
	default:
		return language.Japanese
	}
}

func (l Type) IsSupported() bool {
	_, exist := supported[l]
	return exist
}

var supported = map[Type]struct{}{
	Japanese: {},
}
