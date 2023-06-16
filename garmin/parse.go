package garmin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/mdordoy/peloton-to-garmin/peloton"
	"github.com/pkg/errors"
)

const milesToMetersDistance = 1609.344
const milesPHToMetersPerSecond = 2.237

func ParsePelotonWorkout(workoutDetail peloton.WorkoutDetail, tcxOutToDisk string) (bytes.Buffer, error) {
	tcd := TrainingCenterDatabase{}
	tcd.SchemaLocation = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd"
	tcd.Ns5 = "http://www.garmin.com/xmlschemas/ActivityGoals/v1"
	tcd.Ns3 = "http://www.garmin.com/xmlschemas/ActivityExtension/v2"
	tcd.Ns2 = "http://www.garmin.com/xmlschemas/UserProfile/v2"
	tcd.Xmlns = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"
	tcd.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	tcd.Ns4 = "http://www.garmin.com/xmlschemas/ProfileExtension/v1"

	switch workoutDetail.FitnessDiscipline {
	case "cycling":
		tcd.Activities.Activity.Sport = "Biking"
	case "stretching":
		tcd.Activities.Activity.Sport = "Other"
	default:
		return bytes.Buffer{}, errors.New(fmt.Sprintf("Unsupported sport activity: %s", workoutDetail.FitnessDiscipline))
	}

	startTime := workoutDetail.StartTime.Format("2006-01-02T15:04:05.000Z")
	tcd.Activities.Activity.ID = startTime
	tcd.Activities.Activity.Lap.StartTime = startTime
	tcd.Activities.Activity.Lap.TotalTimeSeconds = workoutDetail.EndTime.Sub(workoutDetail.StartTime).Seconds()
	tcd.Activities.Activity.Lap.DistanceMeters = getDistance(workoutDetail.Summaries)
	tcd.Activities.Activity.Lap.Calories = getTotalCalories(workoutDetail.Summaries)

	summaryMetricData := getSummaryMetricData(workoutDetail.Metrics)

	tcd.Activities.Activity.Lap.AverageHeartRateBpm.Value = summaryMetricData.AvarageHeartRate
	tcd.Activities.Activity.Lap.MaximumHeartRateBpm.Value = summaryMetricData.MaxHeartRate
	tcd.Activities.Activity.Lap.Cadence = summaryMetricData.AverageCadence
	tcd.Activities.Activity.Lap.Extensions.LX.MaxBikeCadence = summaryMetricData.MaxBikeCadence
	tcd.Activities.Activity.Lap.Extensions.LX.MaxWatts = summaryMetricData.MaxWatts
	tcd.Activities.Activity.Lap.Intensity = "Active"
	tcd.Activities.Activity.Lap.TriggerMethod = "Manual"

	tcd.Activities.Activity.Lap.MaximumSpeed = summaryMetricData.MaximumSpeed

	tcd.Activities.Activity.Lap.Extensions.LX.AvgSpeed = getAverageSpeed(workoutDetail.AverageSummaries)
	tcd.Activities.Activity.Lap.Extensions.LX.AvgWatts = getAverageWatts(workoutDetail.AverageSummaries)

	tcd.Activities.Activity.Lap.Track.Trackpoint = parseTrackpointData(&workoutDetail)

	buf, err := writeOutTcxData(tcd, workoutDetail.ID, tcxOutToDisk)
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "failed to write tcx data to file")
	}

	return buf, nil
}

func parseTrackpointData(data *peloton.WorkoutDetail) []Trackpoint {
	trackpoints := []Trackpoint{}
	intervalTime := data.StartTime
	for index := range data.SecondsSincePedalingStart {
		trackpoint := Trackpoint{}
		if index == 0 {
			trackpoint.Time = intervalTime.Format("2006-01-02T15:04:05.000Z")
		} else {
			intervalTime = intervalTime.Add(time.Second * time.Duration(data.DataGranularityInSeconds))
			trackpoint.Time = intervalTime.Format("2006-01-02T15:04:05.000Z")
		}
		for _, data := range data.Metrics {
			switch data.DisplayName {
			case "Output":
				trackpoint.Extensions.TPX.Watts = int(data.Values[index])
			case "Cadence":
				trackpoint.Cadence = int(data.Values[index])
			case "Heart Rate":
				trackpoint.HeartRateBpm.Value = int(data.Values[index])
			case "Speed":
				trackpoint.Extensions.TPX.Speed = data.Values[index] / milesPHToMetersPerSecond
			}
		}
		trackpoints = append(trackpoints, trackpoint)
	}
	return trackpoints
}

func getDistance(summaryData []peloton.WorkoutDetailSummaries) float64 {
	for _, data := range summaryData {
		switch data.DisplayName {
		case "Distance":
			//Convert miles to meteres
			return data.Value * milesToMetersDistance
		}
	}
	return 0
}

func getTotalCalories(summaryData []peloton.WorkoutDetailSummaries) int {
	for _, data := range summaryData {
		switch data.DisplayName {
		case "Calories":
			//Convert miles to meteres
			return int(data.Value)
		}
	}
	return 0
}

func getSummaryMetricData(data []peloton.WorkoutDetailMetrics) metricDetails {
	metricData := metricDetails{}

	for _, metric := range data {
		switch metric.DisplayName {
		case "Speed":
			//Convert Mph to meters per second
			metricData.MaximumSpeed = metric.MaxValue / milesPHToMetersPerSecond
			continue
		case "Heart Rate":
			metricData.MaxHeartRate = int(metric.MaxValue)
			metricData.AvarageHeartRate = int(metric.AverageValue)
			continue
		case "Cadence":
			metricData.MaxBikeCadence = int(metric.MaxValue)
			metricData.AverageCadence = int(metric.AverageValue)
			continue
		case "Output":
			metricData.MaxWatts = int(metric.MaxValue)
		}
	}
	return metricData
}

func getAverageWatts(data []peloton.WorkoutDetailAverageSummaries) int {
	for _, metric := range data {
		switch metric.DisplayName {
		case "Avg Output":
			return int(metric.Value)
		}
	}
	return 0
}

func getAverageSpeed(data []peloton.WorkoutDetailAverageSummaries) float64 {
	for _, metric := range data {
		switch metric.DisplayName {
		case "Avg Speed":
			return metric.Value
		}
	}
	return 0
}

func writeOutTcxData(data TrainingCenterDatabase, workoutID, tcxOutToDisk string) (bytes.Buffer, error) {

	file, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "failed to marshall xml")
	}
	if tcxOutToDisk != "" {
		exists, err := exists(tcxOutToDisk)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, "issue with path provided to write tcx file out to disk")
		}

		if !exists {
			return bytes.Buffer{}, errors.New("Path provided to write tcx to disk does not exist, please fix")
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.tcx", tcxOutToDisk, workoutID), file, 0644)
		if err != nil {
			return bytes.Buffer{}, errors.Wrap(err, "failed to write file")
		}
	}
	var buf bytes.Buffer
	err = xml.NewEncoder(&buf).Encode(data)
	if err != nil {
		return bytes.Buffer{}, errors.Wrap(err, "Failed to encode garmin data for upload")
	}
	return buf, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !stat.IsDir() {
		return false, errors.New("Path provided is not a directory")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, nil
}
