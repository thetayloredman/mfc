package jsonsigning

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thetayloredman/mfc/crypto/ed25519"
)

func TestVerifyJSON(t *testing.T) {
	var keyResponse map[string]any
	err := json.Unmarshal([]byte(`{"old_verify_keys":{},"server_name":"zirco.dev","signatures":{"zirco.dev":{"ed25519:a_mkrB":"oMfvqBYBOooU6u/D3jmqW1JHoJY6e5w+uPFL2FhOdCIRpGcP0+++Nq9xQhYaxwy4UPt3Ox7bgIF4UsbVCWl5AQ"}},"valid_until_ts":1785515203737,"verify_keys":{"ed25519:a_mkrB":{"key":"inYHx5OsCHFuOzZAqh2XYhuWz7FW0RslgVVWLJj9Ajs"}}}`), &keyResponse)

	assert.NoError(t, err)

	fetcher := KeyFetcher(func(serverName, keyID string) (ed25519.PublicKey, error) {
		key := keyResponse["verify_keys"].(map[string]any)[keyID].(map[string]any)["key"].(string)

		decoded, err := base64.RawStdEncoding.DecodeString(key)

		if err != nil {
			return nil, err
		}

		return ed25519.PublicKey(decoded), nil
	})

	err = VerifyJSON(keyResponse, "zirco.dev", fetcher)
	assert.NoError(t, err)
}

func TestSignJSON(t *testing.T) {
	// This uses the Cryptographic Test Vectors from https://spec.matrix.org/v1.19/appendices/#json-signing
	seed, err := base64.RawStdEncoding.DecodeString("YJDBA9Xnr2sVqXD9Vj7XVUnmFZcZrlw8Md7kMW+3XA1")
	assert.NoError(t, err)
	key := ed25519.NewKeyFromSeed(seed)

	signingKey := SigningKey{
		ServerName: "domain",
		KeyID:      "ed25519:1",
		PrivateKey: key,
	}

	cases := []struct {
		input  map[string]any
		signed map[string]any
	}{
		{
			input: map[string]any{},
			signed: map[string]any{
				"signatures": map[string]any{
					"domain": map[string]any{
						"ed25519:1": "K8280/U9SSy9IVtjBuVeLr+HpOB4BQFWbg+UZaADMtTdGYI7Geitb76LTrr5QV/7Xg4ahLwYGYZzuHGZKM5ZAQ",
					},
				},
			},
		},
		{
			input: map[string]any{
				"one": 1,
				"two": "Two",
			},
			signed: map[string]any{
				"one": float64(1),
				"signatures": map[string]any{
					"domain": map[string]any{
						"ed25519:1": "KqmLSbO39/Bzb0QIYE82zqLwsA+PDzYIpIRA2sRQ4sL53+sN6/fpNSoqE7BP7vBZhG6kYdD13EIMJpvhJI+6Bw",
					},
				},
				"two": "Two",
			},
		},
	}

	for _, c := range cases {
		signedJSON, err := SignJSON(c.input, signingKey)
		assert.NoError(t, err)

		var signed map[string]any
		err = json.Unmarshal(signedJSON, &signed)
		assert.NoError(t, err)

		assert.Equal(t, c.signed, signed)
	}
}
