package cmd

import (
	"context"

	"github.com/mdordoy/peloton-to-garmin/logger"
	"github.com/mdordoy/peloton-to-garmin/peloton"
	"github.com/spf13/cobra"
)

var syncConfig struct {
	LogLevel        string
	PrettyLog       bool
	PelotonUsername string
	PelotonPassword string
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Performs the sync",
	RunE:  syncCmd,
}

func syncCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	logger := logger.NewLogger(syncConfig.LogLevel, syncConfig.PrettyLog)
	_ = logger.WithContext(ctx)
	_, err := peloton.NewClient(syncConfig.PelotonUsername, syncConfig.PelotonPassword)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to authenticate with Peloton")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(SyncCmd)
	SyncCmd.Flags().BoolVar(&syncConfig.PrettyLog, "PrettyLogging", false, "Use true for human readable log output")
	SyncCmd.Flags().StringVar(&syncConfig.LogLevel, "loglevel", "info", "Log Level: trace, debug, info, warn,error")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonPassword, "pelotonPassword", "", "peloton Password")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonUsername, "pelotonUsername", "", "peloton Username")
}
