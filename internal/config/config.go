package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env     string        `yaml:"env" env-default:"local"`
	Storage StorageConfig `yaml:"storage" env-required:"true"`
	GRPC    GRPCConfig    `yaml:"grpc_server" env-required:"true"`
	Token   TokenConfig   `yaml:"token" env-required:"true"`
}

type StorageConfig struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
	DBType      string `yaml:"db_type" env-defult:"sqlite"`
}

type GRPCConfig struct {
	Address     string `yaml:"address" env-required:"true"`
	Port        int    `yaml:"port" env-required:"true"`
	Timeout     string `yaml:"timeout" env-required:"true"`
	IdleTimeout string `yaml:"idle_timeout" env-required:"true"`
}
type TokenConfig struct {
	SecretKey      string        `yaml:"secret_key" env-required:"true"`
	ExpirationTime time.Duration `yaml:"expiration_time" env-required:"true"`
}

func MustLoad() *Config {
	// os.Setenv("CONFIG_PATH", "/auth-service/config/config.yaml")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
