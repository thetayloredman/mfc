package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/thetayloredman/mfc/config"
	"github.com/thetayloredman/mfc/crypto/ed25519"
)

func main() {
	if len(os.Args) < 2 {
		panic("Usage: keygen <server_name>")
	}

	serverName := os.Args[1]

	seed := make([]byte, ed25519.SeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	// keyVersion is a 6 byte random string of [a-zA-Z0-9_]
	keyVersion := make([]byte, 6)
	_, err = rand.Read(keyVersion)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(keyVersion); i++ {
		keyVersion[i] = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"[keyVersion[i]%63]
	}

	identity := config.ConfigIdentity{
		ServerName:     serverName,
		KeyId:          "ed25519:" + string(keyVersion),
		PrivateKeySeed: base64.RawStdEncoding.EncodeToString(seed),
	}

	configData := config.Config{
		Identity: identity,
	}

	configText, err := toml.Marshal(configData)
	if err != nil {
		panic(err)
	}

	configPath := config.GetConfigPath()

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("mfc.toml text:\n%s\n", string(configText))
		panic("Config file already exists at " + configPath)
	}

	// create and write the config file
	file, err := os.Create(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(configText)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Config file created at %s\n", configPath)
}
