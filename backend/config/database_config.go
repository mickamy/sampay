package config

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/caarlos0/env/v11"
)

type DatabaseConfigProvider interface {
	DSN() string

	URL() string
}

type EscapableString string

func (s EscapableString) Escape() string {
	return url.QueryEscape(string(s))
}

type DatabaseConfig struct {
	Host       string          `env:"DB_HOST"`
	Port       int             `env:"DB_PORT"`
	Writer     EscapableString `env:"DB_WRITER_USER"`
	WriterPass EscapableString `env:"DB_WRITER_PASSWORD"`
	Reader     EscapableString `env:"DB_READER_USER"`
	ReaderPass EscapableString `env:"DB_READER_PASSWORD"`
	Name       EscapableString `env:"DB_NAME"`
	TimeZone   string          `env:"DB_TIMEZONE"`
	AdminUser  EscapableString `env:"DB_ADMIN_USER" envDefault:"postgres"`
	AdminPass  EscapableString `env:"DB_ADMIN_PASSWORD"`
}

func (c DatabaseConfig) ReaderURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Tokyo",
		c.Reader.Escape(), c.ReaderPass.Escape(), c.Host, c.Port, c.Name.Escape(),
	)
}

func (c DatabaseConfig) WriterURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=Asia/Tokyo",
		c.Writer.Escape(), c.WriterPass.Escape(), c.Host, c.Port, c.Name.Escape(),
	)
}

func (c DatabaseConfig) WriterDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.Writer.Escape(), c.Name.Escape(), c.WriterPass.Escape(), c.TimeZone,
	)
}

func (c DatabaseConfig) ReaderDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.Reader.Escape(), c.Name.Escape(), c.ReaderPass.Escape(), c.TimeZone,
	)
}

func (c DatabaseConfig) AdminDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s sslmode=disable",
		c.Host, c.Port, c.AdminUser.Escape(), c.AdminPass.Escape(), "postgres", c.TimeZone,
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
			database.Port == 0 ||
			database.Writer == "" ||
			database.WriterPass == "" ||
			database.Reader == "" ||
			database.ReaderPass == "" ||
			database.Name == "" ||
			database.TimeZone == "" ||
			database.AdminUser == "" ||
			database.AdminPass == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", database))
		}
	})
	return database
}
