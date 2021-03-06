package cmd

import (
	"github.com/spf13/cobra"

	"github.com/guumaster/hostctl/pkg/host/file"
	"github.com/guumaster/hostctl/pkg/host/render"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [profiles] [flags]",
		Short: "Shows a list of profile names and statuses on your hosts file.",
		Long: `
Shows a list of unique profile names on your hosts file with its status.

The "default" profile is always on and will be skipped.
`,
		RunE: func(cmd *cobra.Command, profiles []string) error {
			src, _ := cmd.Flags().GetString("host-file")

			h, err := file.NewFile(src)
			if err != nil {
				return err
			}

			r := getRenderer(cmd, &render.TableRendererOptions{
				Columns: render.ProfilesOnlyColumns,
			})

			h.ProfileStatus(r, profiles)

			return err
		},
	}
}
