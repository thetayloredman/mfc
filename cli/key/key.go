package sub_key

import (
	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
)

func NewKeyCommand(cfg *config.Config) *cobra.Command {
	keyCmd := &cobra.Command{
		Use:   "key",
		Short: "Remote server key management commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	keyCmd.AddCommand(NewQueryDirectlyCommand(cfg))

	return keyCmd
}
