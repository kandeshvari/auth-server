package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/kandeshvari/auth-server/logger"
	"time"
)

type Server struct {
	Address string `yaml:"address"`
	TLSDir  string `yaml:"tls_dir"`
}

type Db struct {
	Driver    string `yaml:"driver"`
	DbConnect string `yaml:"db_connect"`
}

type Auth struct {
	LoginDelay     time.Duration `yaml:"login_delay"`
	SecretKey      string        `yaml:"secret_key"`
	Timeout        time.Duration `yaml:"timeout"`
	RefreshTimeout time.Duration `yaml:"refresh_timeout"`
}

type Config struct {
	Server   Server `yaml:"server"`
	Database Db     `yaml:"database"`
	Auth     Auth   `yaml:"auth"`

	Logger []logger.LoggerConfig `yaml:"logger"`
}

var defaultConfig = &Config{
	Server: Server{
		Address: "127.0.0.1:9000",
		TLSDir:  "/var/lib/auth-server",
	},
	Database: Db{
		Driver:    "sqlite3",
		DbConnect: "auth.db",
	},
	Auth: Auth{
		LoginDelay:     1000,
		SecretKey:      "",
		Timeout:        30,
		RefreshTimeout: 60 * 24 * 14,
	},

	Logger: []logger.LoggerConfig{
		{Type: "stdout", Name: "default", Format: "[%{level:.1s} %{time:15:04:05}] %{message}", Level: "debug"},
	},
}

func NewConfig(defaults bool) *Config {
	if defaults {
		return defaultConfig
	} else {
		return &Config{}
	}
}

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read the configuration file: %v", err)
	}

	c := NewConfig(true)
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return nil, fmt.Errorf("unable to decode the configuration: %v", err)
	}

	return c, err
}
