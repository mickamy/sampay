// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"mickamy.com/sampay/config"
)

// Injectors from wire.go:

func InitInfras() (Infras, error) {
	databaseConfig := config.Database()
	readWriter, err := provideReadWriter(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	writer, err := provideWriter(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	reader, err := provideReader(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	kvsConfig := config.KVS()
	client, err := provideKVS(kvsConfig)
	if err != nil {
		return Infras{}, err
	}
	diInfras := Infras{
		ReadWriter: readWriter,
		Writer:     writer,
		Reader:     reader,
		KVS:        client,
	}
	return diInfras, nil
}

// wire.go:

func InitConfigs() Configs {
	return NewConfigs()
}
