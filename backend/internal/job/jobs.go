package job

import (
	"fmt"
	"sort"

	"github.com/mickamy/go-sqs-worker/job"

	"github.com/mickamy/sampay/internal/di"
)

type Jobs struct {
	_                  *di.Infra `inject:"param"`
	*ClaimNotification `inject:""`
}

//go:generate go tool stringer -type=Type
type Type int

const (
	first Type = iota
	ClaimNotificationJob
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
	case ClaimNotificationJob:
		return jobs.ClaimNotification, nil
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
