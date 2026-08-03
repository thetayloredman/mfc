package sub_json

import (
	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
)

func NewJsonCommand(cfg *config.Config) *cobra.Command {
	jsonCmd := &cobra.Command{
		Use:   "json",
		Short: "JSON signing and manipulation commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	jsonCmd.AddCommand(NewSignCommand(cfg))

	return jsonCmd
}
