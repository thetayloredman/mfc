package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/thetayloredman/mfc/config"
	"github.com/thetayloredman/mfc/sender"
)

func NewSendCommand(cfg *config.Config) *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send a request to a remote server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	sendCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Send a GET request to a remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmd.Help()
			}

			serverName := args[0]
			path := args[1]

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			s := sender.NewSender(signingKey)
			code, resp, err := s.Get(serverName, path)
			if err != nil {
				return err
			}

			cmd.Printf("Response code: %d\n", code)
			// pretty print the response body as JSON:
			jsonResp, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Printf("Response body:\n%s\n", jsonResp)
			return nil
		},
	})

	sendCmd.AddCommand(&cobra.Command{
		Use:   "post",
		Short: "Send a POST request to a remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			serverName := args[0]
			path := args[1]
			body := args[2]

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			s := sender.NewSender(signingKey)
			code, resp, err := s.Post(serverName, path, []byte(body))
			if err != nil {
				return err
			}

			cmd.Printf("Response code: %d\n", code)
			// pretty print the response body as JSON:
			jsonResp, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Printf("Response body:\n%s\n", jsonResp)
			return nil
		},
	})

	sendCmd.AddCommand(&cobra.Command{
		Use:   "put",
		Short: "Send a PUT request to a remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			serverName := args[0]
			path := args[1]
			body := args[2]

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			s := sender.NewSender(signingKey)
			code, resp, err := s.Put(serverName, path, []byte(body))
			if err != nil {
				return err
			}

			cmd.Printf("Response code: %d\n", code)
			// pretty print the response body as JSON:
			jsonResp, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Printf("Response body:\n%s\n", jsonResp)
			return nil
		},
	})

	sendCmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Send a DELETE request to a remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			serverName := args[0]
			path := args[1]
			body := args[2]

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			s := sender.NewSender(signingKey)
			code, resp, err := s.Delete(serverName, path, []byte(body))
			if err != nil {
				return err
			}

			cmd.Printf("Response code: %d\n", code)
			// pretty print the response body as JSON:
			jsonResp, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Printf("Response body:\n%s\n", jsonResp)
			return nil
		},
	})

	sendCmd.AddCommand(&cobra.Command{
		Use:   "patch",
		Short: "Send a PATCH request to a remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			serverName := args[0]
			path := args[1]
			body := args[2]

			signingKey, err := cfg.AsSigningKey()
			if err != nil {
				return err
			}

			s := sender.NewSender(signingKey)
			code, resp, err := s.Patch(serverName, path, []byte(body))
			if err != nil {
				return err
			}

			cmd.Printf("Response code: %d\n", code)
			// pretty print the response body as JSON:
			jsonResp, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Printf("Response body:\n%s\n", jsonResp)
			return nil
		},
	})

	return sendCmd
}
