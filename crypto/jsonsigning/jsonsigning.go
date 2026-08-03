package jsonsigning

import (
	"encoding/base64"
	"fmt"

	"github.com/thetayloredman/mfc/crypto/canonicaljson"
	"github.com/thetayloredman/mfc/crypto/ed25519"
)

type SigningKey struct {
	ServerName string
	KeyID      string
	PrivateKey ed25519.PrivateKey
}
type KeyFetcher func(serverName, keyID string) (ed25519.PublicKey, error)

func signingObject(item map[string]any) map[string]any {
	m := make(map[string]any, len(item))
	for k, v := range item {
		if k == "signatures" || k == "unsigned" {
			continue
		}
		m[k] = v
	}
	return m
}

func SignJSONGetSignature(item map[string]any, key SigningKey) (string, error) {
	canonicalJSON, err := canonicaljson.Marshal(signingObject(item))
	if err != nil {
		return "", err
	}

	signature := ed25519.Sign(key.PrivateKey, canonicalJSON)

	return base64.RawStdEncoding.EncodeToString(signature), nil
}

func SignJSON(item map[string]any, key SigningKey) ([]byte, error) {
	signature, err := SignJSONGetSignature(item, key)
	if err != nil {
		return nil, err
	}

	item["signatures"] = map[string]map[string]string{
		key.ServerName: {
			key.KeyID: base64.RawStdEncoding.EncodeToString([]byte(signature)),
		},
	}

	return canonicaljson.Marshal(item)
}

func VerifyJSON(item map[string]any, serverName string, keyFetcher KeyFetcher) error {
	if item["signatures"] == nil {
		return fmt.Errorf("no signatures found in JSON")
	}

	if item["signatures"].(map[string]any)[serverName] == nil {
		return fmt.Errorf("no signatures found for server: %s", serverName)
	}

	serverSignatures := item["signatures"].(map[string]any)[serverName].(map[string]any)
	// we should be able to fetch at least one of the keys for this server, otherwise we can't verify the signature
	wasAbleToFindKey := false

	for keyID, signatureB64 := range serverSignatures {
		publicKey, err := keyFetcher(serverName, keyID)
		if err != nil {
			continue
		}
		wasAbleToFindKey = true

		signature, err := base64.RawStdEncoding.DecodeString(signatureB64.(string))
		if err != nil {
			return fmt.Errorf("failed to decode signature for key %s: %w", keyID, err)
		}

		canonicalJSON, err := canonicaljson.Marshal(signingObject(item))

		if err != nil {
			return fmt.Errorf("failed to marshal JSON for verification: %w", err)
		}

		if !ed25519.Verify(publicKey, canonicalJSON, signature) {
			return fmt.Errorf("signature verification failed for key %s", keyID)
		}
	}

	if !wasAbleToFindKey {
		return fmt.Errorf("no valid keys found for server: %s", serverName)
	}

	return nil
}
