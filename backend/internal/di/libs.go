package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/lib/aws/s3"
	"mickamy.com/sampay/internal/lib/aws/ses"
	"mickamy.com/sampay/internal/lib/oauth"
)

type Libs struct {
	OAuthGoogle oauth.Google
	S3          s3.Client
	SES         ses.Client
}

//lint:ignore U1000 used by wire
var libSet = wire.NewSet(
	oauth.NewGoogle,
	s3.New,
	ses.New,
)
