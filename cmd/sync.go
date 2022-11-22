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
	PelotonAPIHost  string
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
	peloClient, err := peloton.NewClient(syncConfig.PelotonUsername, syncConfig.PelotonPassword, syncConfig.PelotonAPIHost)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to authenticate with Peloton")
	}
	workouts, err := peloClient.GetWorkouts()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get users workouts")
	}

	if workouts.Count == 0 {
		logger.Info().Msg("No workouts found")
		return nil
	}

	for _, workout := range workouts.Data {
		workoutDetails, err := peloClient.GetWorkoutDetails(workout.ID, 5)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to get workout with ID %s, skipping", workout.ID)
			continue
		}
		logger.Info().Msgf("workout Details %v", workoutDetails)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(SyncCmd)
	SyncCmd.Flags().BoolVar(&syncConfig.PrettyLog, "PrettyLogging", false, "Use true for human readable log output")
	SyncCmd.Flags().StringVar(&syncConfig.LogLevel, "loglevel", "info", "Log Level: trace, debug, info, warn,error")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonPassword, "pelotonPassword", "", "peloton Password")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonUsername, "pelotonUsername", "", "peloton Username")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonAPIHost, "PelotonAPIHost", "api.onepeloton.com", "The Peloton API host")
}
