package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"applicationDesignTest/internal/consts"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf(consts.EnvVarNotSet, "CONFIG_PATH")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf(consts.CfgNotExists, configPath)
	}

	c, err := getConfig(configPath)
	if err != nil {
		log.Fatalf(consts.CfgErrorGet, err.Error())
	}

	if err = validateConfig(c); err != nil {
		log.Fatalf(consts.CfgIsInvalid, err.Error())
	}

	return c
}

func getConfig(configPath string) (*Config, error) {

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.UnmarshalStrict(configFile, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func validateConfig(c *Config) error {
	if c.Env == "" {
		return fmt.Errorf(consts.CfgNoEnv)
	}

	if !strings.EqualFold(c.Env, consts.EnvLocal) &&
		!strings.EqualFold(c.Env, consts.EnvDevelop) &&
		!strings.EqualFold(c.Env, consts.EnvProduction) {
		return fmt.Errorf(consts.CfgInvalidEnv, c.Env)
	}

	if c.StoragePath == "" {
		return fmt.Errorf(consts.CfgNoStorage)
	}

	if !strings.EqualFold(c.StoragePath, consts.StorageInMemory) {
		if _, err := os.Stat(c.StoragePath); !os.IsNotExist(err) {
			return fmt.Errorf(consts.CfgInvalidStorage, c.StoragePath)
		}
	}

	return nil
}
