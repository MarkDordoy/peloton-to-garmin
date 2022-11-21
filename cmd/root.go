package cmd

import (
	"github.com/mdordoy/peloton-to-garmin/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "Root Command",
	Short: "Root Command",
	Long:  "Root Command",

	Version: version.GetVersion(),

	RunE: rootCmd,
}

func rootCmd(cmd *cobra.Command, args []string) error {
	return errors.New("subcommand required")
}

func init() {
	RootCmd.SetVersionTemplate(`{{printf "%s" .Version}}
`)
}
