package client

import (
	"encoding/json"
	"fmt"

	"github.com/thetayloredman/mfc/client/keyfetcher"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

type OldVerifyKey struct {
	ExpiredTs int64  `json:"expired_ts"`
	Key       string `json:"key"`
}
type VerifyKey struct {
	Key string `json:"key"`
}

type DirectKeyResponse struct {
	OldVerifyKeys map[string]OldVerifyKey `json:"old_verify_keys"`
	ServerName    string                  `json:"server_name"`
	ValidUntilTs  int64                   `json:"valid_until_ts"`
	VerifyKeys    map[string]VerifyKey    `json:"verify_keys"`
}

func (c *Client) QueryForKeyDirectly(serverName string) (DirectKeyResponse, error) {
	fmt.Printf("client: Querying for key directly from server %s\n", serverName)

	respCode, response, err := c.Sender.Get(serverName, "/_matrix/key/v2/server")

	if err != nil {
		return DirectKeyResponse{}, fmt.Errorf("failed to query for key directly from server %s: %v", serverName, err)
	}

	if respCode != 200 {
		return DirectKeyResponse{}, fmt.Errorf("failed to query for key directly from server %s: got response code %d", serverName, respCode)
	}

	var keyResponse DirectKeyResponse
	untypedResponse, err := json.Marshal(response)
	if err != nil {
		return DirectKeyResponse{}, fmt.Errorf("failed to marshal response from server %s: %v", serverName, err)
	}

	err = json.Unmarshal(untypedResponse, &keyResponse)
	if err != nil {
		return DirectKeyResponse{}, fmt.Errorf("failed to unmarshal response from server %s: %v", serverName, err)
	}

	// make a map[string]string of known keys for the ProvidedKeyFetcher
	knownKeys := make(map[string]string)
	for keyID, verifyKey := range keyResponse.VerifyKeys {
		knownKeys[keyID] = verifyKey.Key
	}

	fetcher := keyfetcher.NewProvidedKeyFetcher(keyResponse.ServerName, knownKeys)

	verificationError := jsonsigning.VerifyJSON(response, keyResponse.ServerName, fetcher)

	if verificationError != nil {
		return DirectKeyResponse{}, fmt.Errorf("failed to verify response from server %s: %v", serverName, verificationError)
	}

	return keyResponse, nil
}
