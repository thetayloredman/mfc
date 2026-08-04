package sub_key

import (
	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/client"
	"github.com/thetayloredman/mfc/config"
)

func NewQueryRemoteCommand(cfg *config.Config) *cobra.Command {
	queryDirectlyCmd := &cobra.Command{
		Use:   "query-remote",
		Short: "Query a notary server for another server's keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return cmd.Help()
			}

			notaryServer := args[0]
			targetServer := args[1]
			var specificKey string = ""
			if len(args) == 3 {
				specificKey = args[2]
			}

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			client := client.NewClient(signingKey)
			if specificKey != "" {
				resp, err := client.QueryForOtherKeysSpecific(notaryServer, targetServer, specificKey)

				if err != nil {
					return err
				}

				cmd.Printf("Response from notary server %s for target server %s and key %s:\n\n", notaryServer, targetServer, specificKey)

				for _, keyResp := range resp.ServerKeys {
					cmd.Printf("Server Name: %s\n", keyResp.ServerName)
					cmd.Printf("Valid Until: %d\n", keyResp.ValidUntilTs)
					for keyID, key := range keyResp.VerifyKeys {
						cmd.Printf("Key ID: %s, Key: %s\n", keyID, key.Key)
					}
					for keyID, oldKey := range keyResp.OldVerifyKeys {
						cmd.Printf("Old Key ID: %s, Key: %s, Expired At: %d\n", keyID, oldKey.Key, oldKey.ExpiredTs)
					}
				}

				return nil
			} else {
				resp, err := client.QueryForOtherKeysAll(notaryServer, targetServer)

				if err != nil {
					return err
				}

				cmd.Printf("Response from notary server %s for target server %s:\n\n", notaryServer, targetServer)

				for _, keyResp := range resp.ServerKeys {
					cmd.Printf("Server Name: %s\n", keyResp.ServerName)
					cmd.Printf("Valid Until: %d\n", keyResp.ValidUntilTs)
					for keyID, key := range keyResp.VerifyKeys {
						cmd.Printf("Key ID: %s, Key: %s\n", keyID, key.Key)
					}
					for keyID, oldKey := range keyResp.OldVerifyKeys {
						cmd.Printf("Old Key ID: %s, Key: %s, Expired At: %d\n", keyID, oldKey.Key, oldKey.ExpiredTs)
					}
				}

				return nil
			}
		},
	}

	return queryDirectlyCmd
}
