package config

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
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
	Host       string          `env:"DB_HOST"            validate:"required"`
	Port       int             `env:"DB_PORT"            validate:"required"`
	Writer     EscapableString `env:"DB_WRITER_USER"     validate:"required"`
	WriterPass EscapableString `env:"DB_WRITER_PASSWORD" validate:"required"`
	Reader     EscapableString `env:"DB_READER_USER"     validate:"required"`
	ReaderPass EscapableString `env:"DB_READER_PASSWORD" validate:"required"`
	Name       EscapableString `env:"DB_NAME"            validate:"required"`
	TimeZone   string          `env:"DB_TIMEZONE"        validate:"required"`
	AdminUser  EscapableString `env:"DB_ADMIN_USER"      envDefault:"postgres" validate:"required"`
	AdminPass  EscapableString `env:"DB_ADMIN_PASSWORD"  validate:"required"`
}

func (c DatabaseConfig) ReaderURL() string {
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?TimeZone=%s&sslmode=disable",
		c.Reader.Escape(), c.ReaderPass.Escape(), hostPort, c.Name.Escape(), c.TimeZone,
	)
}

func (c DatabaseConfig) WriterURL() string {
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?TimeZone=%s&sslmode=disable",
		c.Writer.Escape(), c.WriterPass.Escape(), hostPort, c.Name.Escape(), c.TimeZone,
	)
}

func (c DatabaseConfig) AdminURL() string {
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?TimeZone=%s&sslmode=disable",
		c.AdminUser.Escape(), c.AdminPass.Escape(), hostPort, c.Name.Escape(), c.TimeZone,
	)
}

func (c DatabaseConfig) AdminMaintenanceURL() string {
	hostPort := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	return fmt.Sprintf(
		"postgres://%s:%s@%s/postgres?sslmode=disable",
		c.AdminUser.Escape(), c.AdminPass.Escape(), hostPort,
	)
}

func (c DatabaseConfig) WriterDSN() string {
	return buildDSN(
		[][2]string{
			{"host", c.Host},
			{"port", strconv.Itoa(c.Port)},
			{"user", string(c.Writer)},
			{"password", string(c.WriterPass)},
			{"dbname", string(c.Name)},
			{"TimeZone", c.TimeZone},
			{"sslmode", "disable"},
		},
	)
}

func (c DatabaseConfig) ReaderDSN() string {
	return buildDSN(
		[][2]string{
			{"host", c.Host},
			{"port", strconv.Itoa(c.Port)},
			{"user", string(c.Reader)},
			{"password", string(c.ReaderPass)},
			{"dbname", string(c.Name)},
			{"TimeZone", c.TimeZone},
			{"sslmode", "disable"},
		},
	)
}

func (c DatabaseConfig) AdminDSN() string {
	return buildDSN(
		[][2]string{
			{"host", c.Host},
			{"port", strconv.Itoa(c.Port)},
			{"user", string(c.AdminUser)},
			{"password", string(c.AdminPass)},
			{"dbname", "postgres"},
			{"TimeZone", c.TimeZone},
			{"sslmode", "disable"},
		},
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

func buildDSN(params [][2]string) string {
	var sb strings.Builder
	for i, kv := range params {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(kv[0])
		sb.WriteByte('=')
		sb.WriteString(escapeDSNValue(kv[1]))
	}
	return sb.String()
}

func escapeDSNValue(val string) string {
	if val == "" {
		return "''"
	}

	if !strings.ContainsAny(val, " \t\r\n'\\") {
		return val
	}

	var sb strings.Builder
	sb.WriteByte('\'')
	for _, r := range val {
		if r == '\\' || r == '\'' {
			sb.WriteByte('\\')
		}
		sb.WriteRune(r)
	}
	sb.WriteByte('\'')
	return sb.String()
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

		if err := validator.Struct(context.Background(), &database); err != nil {
			panic(fmt.Errorf("invalid database config: %+v", err))
		}
	})
	return database
}
