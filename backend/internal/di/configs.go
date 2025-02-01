package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/config"
)

type Configs struct {
	AWS      config.AWSConfig
	Common   config.CommonConfig
	Database config.DatabaseConfig
	KVS      config.KVSConfig
	OAuth    config.OAuthConfig
}

//lint:ignore U1000 used by wire
var configSet = wire.NewSet(
	config.AWS,
	config.Common,
	config.Database,
	config.KVS,
	config.OAuth,
)
