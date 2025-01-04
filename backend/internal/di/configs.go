package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/config"
)

type Configs struct {
	Common   config.CommonConfig
	Database config.DatabaseConfig
	KVS      config.KVSConfig
}

func NewConfigs() Configs {
	return Configs{
		Common:   config.Common(),
		Database: config.Database(),
		KVS:      config.KVS(),
	}
}

//lint:ignore U1000 used by wire
var configSet = wire.NewSet(
	config.Common,
	config.Database,
	config.KVS,
)
