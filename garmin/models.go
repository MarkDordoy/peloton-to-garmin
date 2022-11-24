package garmin

import (
	"encoding/xml"
)

type TrainingCenterDatabase struct {
	XMLName        xml.Name   `xml:"TrainingCenterDatabase"`
	Text           string     `xml:",chardata"`
	SchemaLocation string     `xml:"xsi:schemaLocation,attr"`
	Ns5            string     `xml:"xmlns:ns5,attr"`
	Ns3            string     `xml:"xmlns:ns3,attr"`
	Ns2            string     `xml:"xmlns:ns2,attr"`
	Xmlns          string     `xml:"xmlns,attr"`
	Xsi            string     `xml:"xmlns:xsi,attr"`
	Ns4            string     `xml:"xmlns:ns4,attr"`
	Activities     Activities `xml:"Activities"`
}

type Activities struct {
	Text     string   `xml:",chardata"`
	Activity Activity `xml:"Activity"`
}

type Activity struct {
	Text  string `xml:",chardata"`
	Sport string `xml:"Sport,attr"`
	ID    string `xml:"Id"`
	Lap   Lap    `xml:"Lap"`
}

type Lap struct {
	Text                string              `xml:",chardata"`
	StartTime           string              `xml:"StartTime,attr"`
	TotalTimeSeconds    float64             `xml:"TotalTimeSeconds"`
	DistanceMeters      float64             `xml:"DistanceMeters"`
	MaximumSpeed        float64             `xml:"MaximumSpeed"`
	Calories            int                 `xml:"Calories"`
	AverageHeartRateBpm AverageHeartRateBpm `xml:"AverageHeartRateBpm"`
	MaximumHeartRateBpm MaximumHeartRateBpm `xml:"MaximumHeartRateBpm"`
	Intensity           string              `xml:"Intensity"`
	Cadence             int                 `xml:"Cadence"`
	TriggerMethod       string              `xml:"TriggerMethod"`
	Track               Track               `xml:"Track"`
	Extensions          Extensions          `xml:"Extensions"`
}

type AverageHeartRateBpm struct {
	Text  string `xml:",chardata"`
	Value int    `xml:"Value"`
}

type MaximumHeartRateBpm struct {
	Text  string `xml:",chardata"`
	Value int    `xml:"Value"`
}

type Extensions struct {
	Text string `xml:",chardata"`
	LX   LX     `xml:"ns3:LX"`
}

type LX struct {
	Text           string  `xml:",chardata"`
	AvgSpeed       float64 `xml:"ns3:AvgSpeed"`
	MaxBikeCadence int     `xml:"ns3:MaxBikeCadence"`
	AvgWatts       int     `xml:"ns3:AvgWatts"`
	MaxWatts       int     `xml:"ns3:MaxWatts"`
}

type TPX struct {
	Text  string  `xml:",chardata"`
	Speed float64 `xml:"ns3:Speed"`
	Watts int     `xml:"ns3:Watts"`
}

type TrackpointExtensions struct {
	Text string `xml:",chardata"`
	TPX  TPX    `xml:"ns3:TPX"`
}

type TrackpointHeartRateBpm struct {
	Text  string `xml:",chardata"`
	Value int    `xml:"Value"`
}

type Track struct {
	Text       string       `xml:",chardata"`
	Trackpoint []Trackpoint `xml:"Trackpoint"`
}

type Trackpoint struct {
	Text         string                 `xml:",chardata"`
	Time         string                 `xml:"Time"`
	HeartRateBpm TrackpointHeartRateBpm `xml:"HeartRateBpm"`
	Cadence      int                    `xml:"Cadence"`
	Extensions   TrackpointExtensions   `xml:"Extensions"`
}

type metricDetails struct {
	MaximumSpeed     float64
	MaxHeartRate     int
	AvarageHeartRate int
	MaxBikeCadence   int
	MaxWatts         int
	AverageCadence   int
}
