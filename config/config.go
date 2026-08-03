package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type ConfigIdentity struct {
	ServerName string `toml:"server_name"`
	KeyId      string `toml:"key_id"`
	PrivateKey string `toml:"private_key"`
}

type Config struct {
	Identity ConfigIdentity `toml:"identity"`
}

func getConfigPath() string {
	if path := os.Getenv("MFC_CONFIG_PATH"); path != "" {
		return path
	}
	return "mfc.toml"
}

func LoadConfig() (*Config, error) {
	configPath := getConfigPath()

	var config Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}

	if config.Identity.ServerName == "" || config.Identity.KeyId == "" || config.Identity.PrivateKey == "" {
		return nil, fmt.Errorf("invalid config: missing required fields in identity")
	}

	return &config, nil
}
