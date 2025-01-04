//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
)

func InitConfigs() Configs {
	return NewConfigs()
}

func InitInfras() (Infras, error) {
	wire.Build(infras, configSet, wire.Struct(new(Infras), "*"))
	return Infras{}, nil
}
