package cmd

import (
	"github.com/mdordoy/peloton-to-garmin/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var syncConfig struct {
	Path string
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Performs the sync",
	RunE:  syncCmd,
}

func syncCmd(cmd *cobra.Command, args []string) error {
	_, err := config.ReadConfig(syncConfig.Path)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading sync configuration at %s", syncConfig.Path)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(SyncCmd)
	SyncCmd.Flags().StringVarP(&syncConfig.Path, "config", "c", "config.yml", "Configuration file. Default is \"config.yml\"")
}
