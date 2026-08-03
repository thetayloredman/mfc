package cli

import (
	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
	"github.com/thetayloredman/mfc/resolver"
)

func NewResolveCommand(cfg *config.Config) *cobra.Command {
	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a server name using the server-server resolution algorithm",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			serverName := args[0]
			resolved, err := resolver.Resolve(serverName)
			if err != nil {
				return err
			}

			cmd.Printf("Resolved destination:\n")
			cmd.Printf("  Endpoint: %s\n", resolved.Endpoint)
			cmd.Printf("  HostHeader: %s\n", resolved.HostHeader)
			cmd.Printf("  ResolutionStep: %s\n", resolved.ResolutionStep)

			return nil
		},
	}

	return resolveCmd
}
