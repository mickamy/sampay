package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type DatabaseConfigProvider interface {
	DSN() string

	URL() string
}

type DatabaseConfig struct {
	Host          string `env:"DB_HOST"`
	Port          string `env:"DB_PORT"`
	User          string `env:"DB_WRITER_USER"`
	Password      string `env:"DB_WRITER_PASSWORD"`
	Name          string `env:"DB_NAME"`
	TimeZone      string `env:"DB_TIMEZONE"`
	AdminUser     string `env:"DB_ADMIN_USER" envDefault:"postgres"`
	AdminPassword string `env:"DB_ADMIN_PASSWORD"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name, c.TimeZone,
	)
}

func (c DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Tokyo",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

func (c DatabaseConfig) AdminDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=postgres password=%s dbname=postgres TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.AdminPassword, c.TimeZone,
	)
}

func (c DatabaseConfig) AdminURL() string {
	return fmt.Sprintf(
		"postgres://postgres:%s@%s:%s/postgres?TimeZone=%s&sslmode=disable",
		c.AdminPassword, c.Host, c.Port, c.TimeZone,
	)
}

var (
	databaseOnce sync.Once
	database     DatabaseConfig
	_            DatabaseConfigProvider = (*DatabaseConfig)(nil)
)

func Database() DatabaseConfig {
	databaseOnce.Do(func() {
		if err := env.Parse(&database); err != nil {
			panic(err)
		}

		if database.Host == "" ||
			database.Port == "" ||
			database.User == "" ||
			database.Password == "" ||
			database.Name == "" ||
			database.TimeZone == "" ||
			database.AdminUser == "" ||
			database.AdminPassword == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", database))
		}
	})
	return database
}
