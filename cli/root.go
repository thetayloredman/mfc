package cli

import (
	"github.com/spf13/cobra"
	sub_json "github.com/thetayloredman/mfc/cli/json"
	"github.com/thetayloredman/mfc/config"
)

func NewRootCommand(cfg *config.Config) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mfc",
		Short: "Matrix server-server API shenanigans",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(sub_json.NewJsonCommand(cfg))

	return rootCmd
}
