package config

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/thetayloredman/mfc/crypto/ed25519"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

type ConfigIdentity struct {
	ServerName     string `toml:"server_name"`
	KeyId          string `toml:"key_id"`
	PrivateKey     string `toml:"private_key"`
	PrivateKeySeed string `toml:"private_key_seed"`
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

	if config.Identity.ServerName == "" || config.Identity.KeyId == "" {
		return nil, fmt.Errorf("invalid config: missing required fields in identity")
	}

	if config.Identity.PrivateKey == "" && config.Identity.PrivateKeySeed == "" {
		return nil, fmt.Errorf("invalid config: either private_key or private_key_seed must be provided in identity")
	}

	if config.Identity.PrivateKey != "" && config.Identity.PrivateKeySeed != "" {
		return nil, fmt.Errorf("invalid config: both private_key and private_key_seed are provided; only one should be set")
	}

	return &config, nil
}

func (cfg *Config) AsSigningKey() (jsonsigning.SigningKey, error) {
	if cfg.Identity.PrivateKeySeed != "" {
		seedBytes, err := base64.RawStdEncoding.DecodeString(cfg.Identity.PrivateKeySeed)
		if err != nil {
			return jsonsigning.SigningKey{}, fmt.Errorf("failed to decode private key seed: %v", err)
		}

		privateKey := ed25519.NewKeyFromSeed(seedBytes)

		key := jsonsigning.SigningKey{
			ServerName: cfg.Identity.ServerName,
			KeyID:      cfg.Identity.KeyId,
			PrivateKey: privateKey,
		}

		return key, nil
	}

	if cfg.Identity.PrivateKey != "" {
		privateKeyBytes, err := base64.RawStdEncoding.DecodeString(cfg.Identity.PrivateKey)
		if err != nil {
			return jsonsigning.SigningKey{}, fmt.Errorf("failed to decode private key: %v", err)
		}

		key := jsonsigning.SigningKey{
			ServerName: cfg.Identity.ServerName,
			KeyID:      cfg.Identity.KeyId,
			PrivateKey: privateKeyBytes,
		}

		return key, nil
	}

	return jsonsigning.SigningKey{}, fmt.Errorf("no private key or seed provided in config")
}
