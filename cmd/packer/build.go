package main

import (
	"get.porter.sh/mixin/packer/pkg/packer"
	"github.com/spf13/cobra"
)

func buildBuildCommand(m *packer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate Dockerfile lines for the bundle invocation image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Build()
		},
	}
	return cmd
}
