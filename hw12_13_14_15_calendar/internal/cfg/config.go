package cfg

import (
	"fmt"
	"os"

	yml "gopkg.in/yaml.v2"
)

const (
	MemoryStorage = "memory"
	SQLStorage    = "sql"
)

//go:generate mockgen -source=./internal/cfg/config.go --destination=./test/mocks/cfg/config.go
type Configurable interface {
	Parse(path string) error
	GetLoggerConf() LoggerConf
	GetDBConf() DatabaseConf
	GetServersConf() ServersConf
	GetAppConf() AppConf
}

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Servers  ServersConf
	App      AppConf
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type AppConf struct {
	StorageType string `yaml:"storageType"`
}

type ServersConf struct {
	HTTP HTTPServerConf `yaml:"http"`
	GRPC GRPCServerConf `yaml:"grpc"`
}

type HTTPServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Parse(path string) error {
	configYml, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %v error: %w", path, err)
	}

	err = yml.Unmarshal(configYml, c)
	if err != nil {
		return fmt.Errorf("can't parse %v: %w", path, err)
	}

	return nil
}

func (c *Config) GetLoggerConf() LoggerConf {
	return c.Logger
}

func (c *Config) GetDBConf() DatabaseConf {
	return c.Database
}

func (c *Config) GetServersConf() ServersConf {
	return c.Servers
}

func (c *Config) GetAppConf() AppConf {
	return c.App
}
