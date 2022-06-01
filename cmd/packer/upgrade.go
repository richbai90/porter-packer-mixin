package main

import (
	"get.porter.sh/mixin/packer/pkg/packer"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *packer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the invoke functionality of this mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
