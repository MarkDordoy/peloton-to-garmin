package garmin

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mdordoy/peloton-to-garmin/peloton"
	"github.com/pkg/errors"
)

const milesToMetersDistance = 1609.344
const milesPHToMetersPerSecond = 2.237

func ParsePelotonWorkout(workoutDetail peloton.WorkoutDetail, dataGranularity int) (TrainingCenterDatabase, error) {
	tcd := TrainingCenterDatabase{}
	tcd.SchemaLocation = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd"
	tcd.Ns5 = "http://www.garmin.com/xmlschemas/ActivityGoals/v1"
	tcd.Ns3 = "http://www.garmin.com/xmlschemas/ActivityExtension/v2"
	tcd.Ns2 = "http://www.garmin.com/xmlschemas/UserProfile/v2"
	tcd.Xmlns = "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2"
	tcd.Xsi = "http://www.w3.org/2001/XMLSchema-instance"
	tcd.Ns4 = "http://www.garmin.com/xmlschemas/ProfileExtension/v1"

	startTime := workoutDetail.StartTime.Format("2006-01-02T15:04:05.000Z")
	tcd.Activities.Activity.ID = startTime
	tcd.Activities.Activity.Lap.StartTime = startTime
	tcd.Activities.Activity.Lap.TotalTimeSeconds = workoutDetail.EndTime.Sub(workoutDetail.StartTime).Seconds()

	var totalDistance float64
	var totalCalories int

	totalDistanceData := workoutDetail.Summaries
	for _, data := range totalDistanceData {
		switch data.DisplayName {
		case "Distance":
			//Convert miles to meteres
			totalDistance = data.Value * milesToMetersDistance
		case "Calories":
			totalCalories = int(data.Value)
		}
	}
	tcd.Activities.Activity.Lap.DistanceMeters = totalDistance
	tcd.Activities.Activity.Lap.Calories = totalCalories

	var maximumSpeed float64
	var maxHeartRate int
	var avarageHeartRate int
	var maxBikeCadence int
	var maxWatts int
	var averageCadence int
	for _, data := range workoutDetail.Metrics {
		switch data.DisplayName {
		case "Speed":
			//Convert Mph to meters per second
			maximumSpeed = data.MaxValue / milesPHToMetersPerSecond
		case "Heart Rate":
			maxHeartRate = int(data.MaxValue)
			avarageHeartRate = int(data.AverageValue)
		case "Cadence":
			maxBikeCadence = int(data.MaxValue)
			averageCadence = int(data.AverageValue)
		case "Output":
			maxWatts = int(data.MaxValue)
		}
	}

	tcd.Activities.Activity.Lap.AverageHeartRateBpm.Value = avarageHeartRate
	tcd.Activities.Activity.Lap.MaximumHeartRateBpm.Value = maxHeartRate
	tcd.Activities.Activity.Lap.Cadence = averageCadence
	tcd.Activities.Activity.Lap.Extensions.LX.MaxBikeCadence = maxBikeCadence
	tcd.Activities.Activity.Lap.Extensions.LX.MaxWatts = maxWatts
	tcd.Activities.Activity.Lap.Intensity = "Active"
	tcd.Activities.Activity.Lap.TriggerMethod = "Manual"

	tcd.Activities.Activity.Lap.MaximumSpeed = maximumSpeed

	switch workoutDetail.FitnessDiscipline {
	case "cycling":
		tcd.Activities.Activity.Sport = "Biking"
	default:
		return tcd, errors.New("Unknown sport activity")
	}

	var avgSpeed float64
	var avgWatts int

	for _, data := range workoutDetail.AverageSummaries {
		switch data.DisplayName {
		case "Avg Output":
			avgWatts = int(data.Value)
		case "Avg Speed":
			avgSpeed = data.Value
		}
	}

	tcd.Activities.Activity.Lap.Extensions.LX.AvgSpeed = avgSpeed
	tcd.Activities.Activity.Lap.Extensions.LX.AvgWatts = avgWatts
	intervalTime := workoutDetail.StartTime

	trackpoints := []Trackpoint{}
	for index, _ := range workoutDetail.SecondsSincePedalingStart {
		trackpoint := Trackpoint{}
		if index == 0 {
			trackpoint.Time = intervalTime.Format("2006-01-02T15:04:05.000Z")
		} else {
			intervalTime = intervalTime.Add(time.Second * time.Duration(dataGranularity))
			trackpoint.Time = intervalTime.Format("2006-01-02T15:04:05.000Z")
		}
		for _, data := range workoutDetail.Metrics {
			switch data.DisplayName {
			case "Output":
				trackpoint.Extensions.TPX.Watts = int(data.Values[index])
			case "Cadence":
				trackpoint.Cadence = int(data.Values[index])
			case "Heart Rate":
				trackpoint.HeartRateBpm.Value = int(data.Values[index])
			case "Speed":
				trackpoint.Extensions.TPX.Speed = data.Values[index]
			}
		}
		trackpoints = append(trackpoints, trackpoint)
	}

	tcd.Activities.Activity.Lap.Track.Trackpoint = trackpoints

	file, err := xml.MarshalIndent(tcd, "", "")
	if err != nil {
		return tcd, errors.Wrap(err, "failed to marshall xml")
	}
	err = ioutil.WriteFile(fmt.Sprintf("/home/mdordoy/github/public/peloton-to-garmin/temp/%s.xml", workoutDetail.Id), file, 0644)
	if err != nil {
		return tcd, errors.Wrap(err, "failed to write file")
	}
	return tcd, nil
}
