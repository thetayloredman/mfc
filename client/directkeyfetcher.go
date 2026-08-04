package client

import (
	"encoding/base64"
	"fmt"

	"github.com/thetayloredman/mfc/crypto/ed25519"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

func NewDirectKeyFetcher(client *Client) jsonsigning.KeyFetcher {
	return func(serverName, keyID string) (ed25519.PublicKey, error) {
		keyResponse, err := client.QueryForKeyDirectly(serverName)
		if err != nil {
			return nil, err
		}

		key := keyResponse.VerifyKeys[keyID].Key
		if key == "" {
			return nil, fmt.Errorf("don't have a key for server %s with key ID %s", serverName, keyID)
		}

		decoded, err := base64.RawStdEncoding.DecodeString(key)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 key for server %s with key ID %s: %v", serverName, keyID, err)
		}
		fmt.Printf("keyfetcher/direct: Returning direct key for server %s with key ID %s: %x\n", serverName, keyID, decoded)

		return ed25519.PublicKey(decoded), nil
	}
}
