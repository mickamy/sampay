package job

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mickamy/go-sqs-worker/job"

	"mickamy.com/sampay/internal/lib/validator"
)

type Jobs struct {
	SendEmail
}

//go:generate stringer -type=Type
type Type int

const (
	first Type = iota
	SendEmailJob
	last
)

func Get(s string, jobs Jobs) (job.Job, error) {
	types := _types()

	idx := sort.Search(len(types), func(i int) bool {
		return types[i].String() >= s
	})

	if idx == len(types) || types[idx].String() != s {
		return nil, fmt.Errorf("unknown job type: [%s]", s)
	}

	switch types[idx] {
	case first, last:
		return nil, fmt.Errorf("type `first` and `last` should not be used")
	case SendEmailJob:
		return jobs.SendEmail, nil
	}

	return nil, fmt.Errorf("unknown job type: %s", s)
}

func _types() []Type {
	var types []Type
	for i := first + 1; i < last; i++ {
		types = append(types, i)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].String() < types[j].String()
	})
	return types
}

func parsePayload(ctx context.Context, payloadStr string, to any) error {
	if err := json.Unmarshal([]byte(payloadStr), &to); err != nil {
		return fmt.Errorf("payload unmarshalling failed: %w", err)
	}
	if err := validator.Struct(ctx, to); err != nil {
		return fmt.Errorf("payload validation failed: %s", err)
	}

	return nil
}
