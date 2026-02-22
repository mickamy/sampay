package database

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray implements sql.Scanner and driver.Valuer for PostgreSQL TEXT[] columns.
type StringArray []string

func (a *StringArray) Scan(src any) error {
	if src == nil {
		*a = nil
		return nil
	}
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("StringArray.Scan: expected string, got %T", src)
	}
	*a = parsePostgresArray(s)
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil //nolint:nilnil // nil is the correct driver.Value for a NULL SQL column
	}
	elems := make([]string, len(a))
	for i, s := range a {
		escaped := strings.ReplaceAll(s, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		elems[i] = `"` + escaped + `"`
	}
	return "{" + strings.Join(elems, ",") + "}", nil
}

func parsePostgresArray(s string) []string {
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return nil
	}
	inner := s[1 : len(s)-1]
	if inner == "" {
		return []string{}
	}
	var result []string
	var current strings.Builder
	inQuote := false
	escaped := false
	for i := range len(inner) {
		c := inner[i]
		if escaped {
			current.WriteByte(c)
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '"' {
			inQuote = !inQuote
			continue
		}
		if c == ',' && !inQuote {
			result = append(result, current.String())
			current.Reset()
			continue
		}
		current.WriteByte(c)
	}
	result = append(result, current.String())
	return result
}
