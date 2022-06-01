package main

import (
	"get.porter.sh/mixin/packer/pkg/packer"
	"github.com/spf13/cobra"
)

var (
	commandFile string
)

func buildInstallCommand(m *packer.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Execute the install functionality of this mixin",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			//Do something here if needed
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Execute()
		},
	}
	return cmd
}
