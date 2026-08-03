package sub_json

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
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

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return fmt.Errorf("failed to get signing key: %v", err)
			}

			signedJSON, err := jsonsigning.SignJSON(jsonData, signingKey)
			if err != nil {
				return fmt.Errorf("failed to sign JSON: %v", err)
			}
			fmt.Println(string(signedJSON))
			return nil
		},
	}

	return signCmd
}
