package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/lib/aws/s3"
)

type Libs struct {
	s3.Client
}

//lint:ignore U1000 used by wire
var libSet = wire.NewSet(
	s3.New,
)
