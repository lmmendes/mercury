package config

import (
	"mercury/internal/logger"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

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
	Database struct {
		Driver string `koanf:"driver"`
		URL    string `koanf:"url"`
	} `koanf:"database"`
	Logging struct {
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
   		port: ":1025"
	    domain: "localhost"
database:
	driver: "sqlite3"
	url: "./email.db"
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
