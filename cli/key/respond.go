package sub_key

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/client"
	"github.com/thetayloredman/mfc/config"
	"github.com/thetayloredman/mfc/crypto/jsonsigning"
)

func NewRespondCommand(cfg *config.Config) *cobra.Command {
	respondCmd := &cobra.Command{
		Use:   "respond",
		Short: "Generate a /_matrix/key/v2/server response",
		RunE: func(cmd *cobra.Command, args []string) error {
			// we issue key responses valid for 1000 years. hopefully by
			// then we have a better protocol.

			currentTime := time.Now().Unix()
			validUntil := currentTime + 1000*365*24*60*60

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}
			publicKey := signingKey.PrivateKey.Public()
			base64Pubkey := base64.RawStdEncoding.EncodeToString(publicKey.([]byte))

			verifyKeys := map[string]client.VerifyKey{}
			verifyKeys[signingKey.KeyID] = client.VerifyKey{
				Key: base64Pubkey,
			}

			keyResponse := client.KeyResponse{
				ServerName:    signingKey.ServerName,
				VerifyKeys:    verifyKeys,
				OldVerifyKeys: map[string]client.OldVerifyKey{},
				ValidUntilTs:  validUntil * 1000, // milliseconds
			}

			// convert into map[string]any using json for signing
			jsonResponse, err := json.Marshal(keyResponse)
			if err != nil {
				return err
			}

			var jsonMap map[string]any
			err = json.Unmarshal(jsonResponse, &jsonMap)
			if err != nil {
				return err
			}

			signedResponse, err := jsonsigning.SignJSON(jsonMap, signingKey)
			if err != nil {
				return err
			}

			cmd.Println(string(signedResponse))
			return nil
		},
	}

	return respondCmd
}
