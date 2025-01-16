package job

import (
	"fmt"
	"sort"

	"github.com/mickamy/go-sqs-worker/job"
)

type Jobs struct {
	SendEmail
}

//go:generate stringer -type=jobType
type Type int

const (
	first Type = iota
	SendEmailType
	last
)

func Get(s string, jobs Jobs) (job.Job, error) {
	types := types()

	idx := sort.Search(len(types), func(i int) bool {
		return types[i].String() >= s
	})

	if idx == len(types) || types[idx].String() != s {
		return nil, fmt.Errorf("unknown job type: [%s]", s)
	}

	switch types[idx] {
	case first, last:
		return nil, fmt.Errorf("type `first` and `last` should not be used")
	case SendEmailType:
		return jobs.SendEmail, nil
	}

	return nil, fmt.Errorf("unknown job type: %s", s)
}

func types() []Type {
	var types []Type
	for i := first + 1; i < last; i++ {
		types = append(types, i)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].String() < types[j].String()
	})
	return types
}
