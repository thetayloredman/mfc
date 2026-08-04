package sender

import (
	"encoding/json"
	"fmt"

	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

type OutboundRequestSignatureContent struct {
	Method      string         `json:"method"`
	Uri         string         `json:"uri"`
	Origin      string         `json:"origin"`
	Destination string         `json:"destination"`
	Content     map[string]any `json:"content,omitempty"`
}

// Generate the "X-Matrix" authorization header for an outbound federation request.
func (s *Sender) authorizeOutboundRequest(method, uri, destination string, body []byte) (string, error) {
	reqContent := OutboundRequestSignatureContent{
		Method:      method,
		Uri:         uri,
		Origin:      s.identity.ServerName,
		Destination: destination,
	}

	if len(body) > 0 {
		var content map[string]any
		err := json.Unmarshal(body, &content)
		if err != nil {
			return "", fmt.Errorf("failed to parse request body as JSON: %v", err)
		}
		reqContent.Content = content
	}

	// convert to a map[string]any for signing
	reqContentMap := map[string]any{
		"method":      reqContent.Method,
		"uri":         reqContent.Uri,
		"origin":      reqContent.Origin,
		"destination": reqContent.Destination,
	}
	if reqContent.Content != nil {
		reqContentMap["content"] = reqContent.Content
	}

	fmt.Printf("sender/auth-out: Signing request content: %+v\n", reqContentMap)

	signature, err := jsonsigning.SignJSONGetSignature(reqContentMap, s.identity)
	if err != nil {
		return "", fmt.Errorf("failed to sign request content: %v", err)
	}

	fmt.Printf("sender/auth-out: Generated signature: %s\n", signature)

	return "X-Matrix origin=\"" + s.identity.ServerName + "\",destination=\"" + destination + "\",key=\"" + s.identity.KeyID + "\",sig=\"" + signature + "\"", nil
}
