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
	Writer        string `env:"DB_WRITER_USER"`
	WriterPass    string `env:"DB_WRITER_PASSWORD"`
	Reader        string `env:"DB_READER_USER"`
	ReaderPass    string `env:"DB_READER_PASSWORD"`
	Name          string `env:"DB_NAME"`
	TimeZone      string `env:"DB_TIMEZONE"`
	AdminUser     string `env:"DB_ADMIN_USER" envDefault:"postgres"`
	AdminPassword string `env:"DB_ADMIN_PASSWORD"`
}

func (c DatabaseConfig) ReaderURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Tokyo",
		c.Reader, c.ReaderPass, c.Host, c.Port, c.Name,
	)
}

func (c DatabaseConfig) WriterURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Tokyo",
		c.Writer, c.WriterPass, c.Host, c.Port, c.Name,
	)
}

func (c DatabaseConfig) WriterDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.Writer, c.Name, c.WriterPass, c.TimeZone,
	)
}

func (c DatabaseConfig) ReaderDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.Reader, c.Name, c.ReaderPass, c.TimeZone,
	)
}

func (c DatabaseConfig) AdminDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=postgres password=%s dbname=postgres TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.AdminPassword, c.TimeZone,
	)
}

type writerProvider struct {
	cfg DatabaseConfig
}

func (w writerProvider) DSN() string {
	return w.cfg.WriterDSN()
}

func (w writerProvider) URL() string {
	return w.cfg.WriterURL()
}

func (c DatabaseConfig) WriterProvider() DatabaseConfigProvider {
	return writerProvider{cfg: c}
}

type readerProvider struct {
	cfg DatabaseConfig
}

func (r readerProvider) DSN() string {
	return r.cfg.ReaderDSN()
}

func (r readerProvider) URL() string {
	return r.cfg.ReaderURL()
}

func (c DatabaseConfig) ReaderProvider() DatabaseConfigProvider {
	return readerProvider{cfg: c}
}

var (
	databaseOnce sync.Once
	database     DatabaseConfig
)

func Database() DatabaseConfig {
	databaseOnce.Do(func() {
		if err := env.Parse(&database); err != nil {
			panic(err)
		}

		if database.Host == "" ||
			database.Port == "" ||
			database.Writer == "" ||
			database.WriterPass == "" ||
			database.Reader == "" ||
			database.ReaderPass == "" ||
			database.Name == "" ||
			database.TimeZone == "" ||
			database.AdminUser == "" ||
			database.AdminPassword == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", database))
		}
	})
	return database
}
