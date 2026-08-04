package client

import (
	"encoding/base64"
	"fmt"

	"github.com/thetayloredman/mfc/crypto/ed25519"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

func NewProvidedKeyFetcher(knownServerName string, knownServerKeys map[string]string) jsonsigning.KeyFetcher {
	return func(serverName, keyID string) (ed25519.PublicKey, error) {
		if serverName != knownServerName {
			return nil, fmt.Errorf("wasn't expecting to need a key for server %s, only know about %s", serverName, knownServerName)
		}

		key := knownServerKeys[keyID]
		if key == "" {
			return nil, fmt.Errorf("wasn't expecting to need a key for server %s with key ID %s, only know about keys: %v", serverName, keyID, knownServerKeys)
		}

		decoded, err := base64.RawStdEncoding.DecodeString(key)

		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 key for server %s with key ID %s: %v", serverName, keyID, err)
		}
		fmt.Printf("keyfetcher/provided: Returning provided key for server %s with key ID %s: %x\n", serverName, keyID, decoded)

		return ed25519.PublicKey(decoded), nil
	}
}
