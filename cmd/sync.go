package cmd

import (
	"context"
	"strings"

	connect "github.com/abrander/garmin-connect"
	"github.com/mdordoy/peloton-to-garmin/garmin"
	"github.com/mdordoy/peloton-to-garmin/logger"
	"github.com/mdordoy/peloton-to-garmin/peloton"
	"github.com/spf13/cobra"
)

var syncConfig struct {
	LogLevel                string
	PrettyLog               bool
	PelotonUsername         string
	PelotonPassword         string
	PelotonAPIHost          string
	DataGranularity         int
	PelotonWorkoutInstances int
	GarminEmail             string
	GarminPassword          string
	OutTCXFilePath          string
}

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Performs workout syncs from Peloton to Garmin Connect",
	RunE:  syncCmd,
}

func syncCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	logger := logger.NewLogger(syncConfig.LogLevel, syncConfig.PrettyLog)
	_ = logger.WithContext(ctx)

	if syncConfig.GarminEmail == "" {
		logger.Fatal().Msg("Garmin email not provided, this is required")
	}
	if syncConfig.GarminPassword == "" {
		logger.Fatal().Msg("Garmin password not provided, this is required")
	}
	if syncConfig.PelotonUsername == "" {
		logger.Fatal().Msg("Peloton username not provided, this is required")
	}
	if syncConfig.PelotonPassword == "" {
		logger.Fatal().Msg("Peloton password not provided, this is required")
	}
	peloClient, err := peloton.NewClient(syncConfig.PelotonUsername, syncConfig.PelotonPassword, syncConfig.PelotonAPIHost)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to authenticate with Peloton")
	}
	workouts, err := peloClient.GetWorkouts(syncConfig.PelotonWorkoutInstances)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get users workouts")
	}

	if len(workouts) == 0 {
		logger.Info().Msg("No workouts found")
		return nil
	}

	workoutList := []peloton.WorkoutDetail{}

	for _, workout := range workouts {
		workoutDetails, err := peloClient.GetWorkoutDetails(workout, syncConfig.DataGranularity)
		if err != nil {
			logger.Error().Err(err).Msgf("Failed to get workout with ID %s, skipping", workout.ID)
			continue
		}
		logger.Info().Str("Title", workoutDetails.Title).Str("Workout ID", workoutDetails.ID).Str("Workout Date", workoutDetails.StartTime.Format("Mon Jan 2 2006 15:04:05")).Msg("Found Peloton Workout")
		workoutList = append(workoutList, workoutDetails)
	}

	garminClient := garmin.NewClient(syncConfig.GarminEmail, syncConfig.GarminPassword, logger)
	for _, workoutDetail := range workoutList {
		buf, err := garmin.ParsePelotonWorkout(workoutDetail, syncConfig.OutTCXFilePath)
		if err != nil {
			logger.Error().Err(err).Str("Title", workoutDetail.Title).Str("Workout ID", workoutDetail.ID).Msg("Failed to convert peloton data to garmin data")
			continue
		}

		status, err := garminClient.ImportActivity(&buf, connect.ActivityFormatTCX)
		rLogger := logger.With().Str("Title", workoutDetail.Title).Str("Workout ID", workoutDetail.ID).Str("Workout Date", workoutDetail.StartTime.Format("Mon Jan 2 2006 15:04:05")).Logger()
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "Duplicate Activity"):
				rLogger.Info().Msg("Workout already uploaded to garmin")
			case strings.Contains(err.Error(), "202: Accepted"):
				rLogger.Info().Msgf("Workout uploaded to garmin")
			default:
				rLogger.Error().Err(err).Msg("Failed to upload activity to garmin")
			}
			continue
		}

		err = garminClient.RenameActivity(status, workoutDetail.Title)
		if err != nil {
			rLogger.Warn().Err(err).Msgf("Workout uploaded but failed to rename garmin activity to %s", workoutDetail.Title)
			continue
		}
		rLogger.Info().Msgf("Workout uploaded and renamed to %s", workoutDetail.Title)

	}
	logger.Info().Msg("Peloton to Garmin Sync completed")

	return nil
}

func init() {
	RootCmd.AddCommand(SyncCmd)
	SyncCmd.Flags().BoolVar(&syncConfig.PrettyLog, "PrettyLogging", true, "Use true for human readable log output")
	SyncCmd.Flags().StringVar(&syncConfig.LogLevel, "loglevel", "info", "Log Level: trace, debug, info, warn,error")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonPassword, "pelotonPassword", "", "peloton Password")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonUsername, "pelotonUsername", "", "peloton Username")
	SyncCmd.Flags().StringVar(&syncConfig.PelotonAPIHost, "PelotonAPIHost", "api.onepeloton.com", "The Peloton API host")
	SyncCmd.Flags().IntVar(&syncConfig.DataGranularity, "granularity", 1, "Data granularity from Peloton, default every 1 second")
	SyncCmd.Flags().IntVar(&syncConfig.PelotonWorkoutInstances, "workoutCount", 30, "Number of previous workouts you want to pull from Peloton")
	SyncCmd.Flags().StringVar(&syncConfig.GarminPassword, "garminPassword", "", "Garmin Password")
	SyncCmd.Flags().StringVar(&syncConfig.GarminEmail, "garminEmail", "", "Garmin Email")
	SyncCmd.Flags().StringVar(&syncConfig.OutTCXFilePath, "writeTCXToDisk", "", "If you provide an absolute path, the cli will write the tcx file out to disk")
}
