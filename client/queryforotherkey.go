package client

import (
	"encoding/json"
	"fmt"

	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

type SpecificKeyRequests map[string]struct{}

type NotaryKeyRequest struct {
	ServerKeys map[string]SpecificKeyRequests `json:"server_keys"`
}

type NotaryKeyResponse struct {
	ServerKeys []KeyResponse `json:"server_keys"`
}

func (c *Client) QueryForOtherKeys(notaryServer string, request NotaryKeyRequest) (NotaryKeyResponse, error) {
	fmt.Printf("client: Querying %s for keys\n", notaryServer)

	var requestBody []byte
	var err error

	if requestBody, err = json.Marshal(request); err != nil {
		return NotaryKeyResponse{}, fmt.Errorf("failed to marshal request for other keys: %v", err)
	}

	respCode, response, err := c.Sender.Post(notaryServer, "/_matrix/key/v2/query", requestBody)

	if err != nil {
		return NotaryKeyResponse{}, fmt.Errorf("failed to query for other keys from notary: %v", err)
	}

	if respCode != 200 {
		return NotaryKeyResponse{}, fmt.Errorf("failed to query for other keys from notary: got response code %d", respCode)
	}

	var keyResponse NotaryKeyResponse
	untypedResponse, err := json.Marshal(response)
	if err != nil {
		return NotaryKeyResponse{}, fmt.Errorf("failed to marshal response from notary: %v", err)
	}

	err = json.Unmarshal(untypedResponse, &keyResponse)
	if err != nil {
		return NotaryKeyResponse{}, fmt.Errorf("failed to unmarshal response from notary: %v", err)
	}

	// for every response, verify the notary has signed it (using a directkeyfetcher instance)
	for i, keyResp := range keyResponse.ServerKeys {
		fetcher := NewDirectKeyFetcher(c)

		untypedResponseEntry := response["server_keys"].([]any)[i].(map[string]any)

		verificationError := jsonsigning.VerifyJSON(untypedResponseEntry, notaryServer, fetcher)

		if verificationError != nil {
			return NotaryKeyResponse{}, fmt.Errorf("failed to verify response from notary: %v", verificationError)
		}

		// and also verify the server has signed it (using a providedkeyfetcher instance)
		knownKeys := make(map[string]string)
		for keyID, verifyKey := range keyResp.VerifyKeys {
			knownKeys[keyID] = verifyKey.Key
		}

		fetcher = NewProvidedKeyFetcher(keyResp.ServerName, knownKeys)

		verificationError = jsonsigning.VerifyJSON(untypedResponseEntry, keyResp.ServerName, fetcher)

		if verificationError != nil {
			return NotaryKeyResponse{}, fmt.Errorf("failed to verify response from server %s: %v", keyResp.ServerName, verificationError)
		}
	}

	return keyResponse, nil
}

func (c *Client) QueryForOtherKeysAll(notaryServer string, desiredServer string) (NotaryKeyResponse, error) {
	request := NotaryKeyRequest{
		ServerKeys: map[string]SpecificKeyRequests{
			desiredServer: {},
		},
	}

	return c.QueryForOtherKeys(notaryServer, request)
}

func (c *Client) QueryForOtherKeysSpecific(notaryServer string, desiredServer string, desiredKey string) (NotaryKeyResponse, error) {
	request := NotaryKeyRequest{
		ServerKeys: map[string]SpecificKeyRequests{
			desiredServer: {
				desiredKey: {},
			},
		},
	}

	return c.QueryForOtherKeys(notaryServer, request)
}
