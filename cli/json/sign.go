package sub_json

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
	"github.com/thetayloredman/mfc/crypto/ed25519"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

func NewSignCommand(cfg *config.Config) *cobra.Command {
	signCmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign the provided JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			jsonText := args[0]

			var jsonData map[string]any
			err := json.Unmarshal([]byte(jsonText), &jsonData)
			if err != nil {
				return fmt.Errorf("failed to parse JSON: %v", err)
			}

			privateKeyBytes, err := base64.RawStdEncoding.DecodeString(cfg.Identity.PrivateKey)
			if err != nil {
				return fmt.Errorf("failed to decode private key: %v", err)
			}

			privateKey := ed25519.PrivateKey(privateKeyBytes)

			key := jsonsigning.SigningKey{
				ServerName: cfg.Identity.ServerName,
				KeyID:      cfg.Identity.KeyId,
				PrivateKey: privateKey,
			}

			signedJSON, err := jsonsigning.SignJSON(jsonData, key)
			if err != nil {
				return fmt.Errorf("failed to sign JSON: %v", err)
			}
			fmt.Println(string(signedJSON))
			return nil
		},
	}

	return signCmd
}
