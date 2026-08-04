package client

import (
	"fmt"

	"github.com/thetayloredman/mfc/crypto/ed25519"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

// Creates a keyfetcher that attempts the provided keyfetcher list in order, returning the first key it finds for a given server and key ID. If none of the provided keyfetchers return a key, it returns an error.
func NewKeyFetcherCombinator(keyFetchers ...jsonsigning.KeyFetcher) jsonsigning.KeyFetcher {
	return func(serverName, keyID string) (ed25519.PublicKey, error) {
		for _, keyFetcher := range keyFetchers {
			key, err := keyFetcher(serverName, keyID)
			if err == nil {
				return key, nil
			}
		}

		return nil, fmt.Errorf("no key found for server %s with key ID %s from any of the provided key fetchers", serverName, keyID)
	}
}
