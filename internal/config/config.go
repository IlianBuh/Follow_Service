package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env          string  `yaml:"env" env-default:"prod"`
	StorageURL   string  `yaml:"storage-url" env-required:"true"`
	GRPC         GRPCObj `yaml:"grpc"`
	UserInfoPort int     `yaml:"user-info-port" env-required:"true"`
}

type GRPCObj struct {
	Port       int           `yaml:"port" env-default:"20202"`
	Timeout    time.Duration `yaml:"timeout" env-default:"5s"`
	RetryCount int           `yaml:"retry-count" env-default:"5"`
}

const (
	defaultConfigPath = "./config/config.yml"
)

// MustLoad returns new config object. Panics if error occurred.
func MustLoad() *Config {
	path := fetchConfigPath()

	return MustLoadByPath(path)
}

// MustLoadByPath returns new config object by the 'path'. Panics if error occurred.
func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not found: " + path)
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to parse config file" + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches the config file path either from flag "config" or
// environment variable "CONFIG_PATH". If no both of them, default path is returned
//
// flag > env > default
func fetchConfigPath() string {
	res := ""

	flag.StringVar(&res, "config", "", "path to the config file")
	flag.Parse()
	if res != "" {
		return res
	}

	res = os.Getenv("CONFIG_PATH")
	if res != "" {
		return res
	}

	res = defaultConfigPath
	return res
}
