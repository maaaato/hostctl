package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host"
	"github.com/guumaster/hostctl/pkg/host/docker"
	"github.com/guumaster/hostctl/pkg/host/file"
)

func newSyncDockerCmd(removeCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "docker [profile] [flags]",
		Short: "Sync your Docker containers IPs with a profile.",
		Long: `
Reads from Docker the list of containers and add names and IPs to a profile in your hosts file.
`,
		Args: commonCheckArgs,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")
			domain, _ := cmd.Flags().GetString("domain")
			network, _ := cmd.Flags().GetString("network")

			ctx := context.Background()

			p, err := docker.NewProfileFromDocker(ctx, &docker.Options{
				Domain:  domain,
				Network: network,
				Cli:     nil,
			})
			if err != nil {
				return err
			}

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			p.Name = profiles[0]
			p.Status = host.Enabled

			err = h.AddProfile(p)
			if err != nil {
				return err
			}

			return h.Flush()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return postActionCmd(cmd, args, removeCmd, true)
		},
	}
}
