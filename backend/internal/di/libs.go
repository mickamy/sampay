package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/lib/aws/s3"
	"mickamy.com/sampay/internal/lib/aws/ses"
)

type Libs struct {
	S3  s3.Client
	SES ses.Client
}

//lint:ignore U1000 used by wire
var libSet = wire.NewSet(
	s3.New,
	ses.New,
)
