package config

import (
	"mercury/internal/logger"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type DatabaseConfig struct {
	URL             string        `koanf:"url"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

type Config struct {
	Server struct {
		HTTP struct {
			Port string `koanf:"port"`
		} `koanf:"http"`
		SMTP struct {
			Port     string `koanf:"port"`
			Hostname string `koanf:"hostname"`
			Username string `koanf:"username"`
			Password string `koanf:"password"`
		} `koanf:"smtp"`
		IMAP struct {
			Port     string `koanf:"port"`
			Hostname string `koanf:"hostname"`
		}
	} `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Logging  struct {
		Level  logger.Level `koanf:"level"`
		Format string       `koanf:"format"`
	} `koanf:"logging"`
}

func Load(configFile string) (*Config, error) {
	k := koanf.New(".")

	// Load default configuration
	defaultConfig := []byte(`
server:
  http:
    port: ":8080"
  smtp:
    port: ":1025"
    hostname: "localhost"
    username: ""
    password: ""
  imap:
    port: ":1143"
    hostname: "localhost"
database:
  url: "postgres://mercury:mercury@localhost:5432/mercury?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m
logging:
	level: "info"
	format: "json"
`)

	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		// If config file doesn't exist, use only defaults and environment variables
		if err := k.Load(file.Provider(string(defaultConfig)), yaml.Parser()); err != nil {
			return nil, err
		}
	}

	// Load environment variables
	// MERCURY_SERVER_HTTP_PORT, MERCURY_SERVER_SMTP_PORT, etc.
	if err := k.Load(env.Provider("MERCURY_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MERCURY_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
