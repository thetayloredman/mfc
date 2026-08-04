package sub_key

import (
	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/client"
	"github.com/thetayloredman/mfc/config"
)

func NewQueryDirectlyCommand(cfg *config.Config) *cobra.Command {
	queryDirectlyCmd := &cobra.Command{
		Use:   "query-directly",
		Short: "Query a remote server for its public keys directly",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			client := client.NewClient(signingKey)

			serverName := args[0]
			keys, err := client.QueryForKeyDirectly(serverName)
			if err != nil {
				return err
			}

			cmd.Printf("Keys for server %s:\n", serverName)
			for keyID, key := range keys.VerifyKeys {
				cmd.Printf("Key ID: %s, Key: %s\n", keyID, key.Key)
			}
			for keyID, oldKey := range keys.OldVerifyKeys {
				cmd.Printf("Old Key ID: %s, Key: %s, Expired At: %d\n", keyID, oldKey.Key, oldKey.ExpiredTs)
			}
			cmd.Printf("Valid Until: %d\n", keys.ValidUntilTs)

			return nil
		},
	}

	return queryDirectlyCmd
}
