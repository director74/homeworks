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

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Server   ServerConf
	App      AppConf
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type AppConf struct {
	StorageType string `yaml:"storageType"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DatabaseConf struct {
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
}

func NewConfig() Config {
	return Config{}
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
