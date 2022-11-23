package garmin

import "encoding/xml"

type TrainingCenterDatabase struct {
	XMLName        xml.Name `xml:"TrainingCenterDatabase"`
	Text           string   `xml:",chardata"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Ns5            string   `xml:"xmlns:ns5,attr"`
	Ns3            string   `xml:"xmlns:ns3,attr"`
	Ns2            string   `xml:"xmlns:ns2,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xmlns:xsi,attr"`
	Ns4            string   `xml:"xmlns:ns4,attr"`
	Activities     struct {
		Text     string `xml:",chardata"`
		Activity struct {
			Text  string `xml:",chardata"`
			Sport string `xml:"Sport,attr"`
			ID    string `xml:"Id"`
			Lap   struct {
				Text                string  `xml:",chardata"`
				StartTime           string  `xml:"StartTime,attr"`
				TotalTimeSeconds    float64 `xml:"TotalTimeSeconds"`
				DistanceMeters      float64 `xml:"DistanceMeters"`
				MaximumSpeed        float64 `xml:"MaximumSpeed"`
				Calories            int     `xml:"Calories"`
				AverageHeartRateBpm struct {
					Text  string `xml:",chardata"`
					Value int    `xml:"Value"`
				} `xml:"AverageHeartRateBpm"`
				MaximumHeartRateBpm struct {
					Text  string `xml:",chardata"`
					Value int    `xml:"Value"`
				} `xml:"MaximumHeartRateBpm"`
				Intensity     string `xml:"Intensity"`
				Cadence       int    `xml:"Cadence"`
				TriggerMethod string `xml:"TriggerMethod"`
				Track         struct {
					Text       string       `xml:",chardata"`
					Trackpoint []Trackpoint `xml:"Trackpoint"`
				} `xml:"Track"`
				Extensions struct {
					Text string `xml:",chardata"`
					LX   struct {
						Text           string  `xml:",chardata"`
						AvgSpeed       float64 `xml:"ns3:AvgSpeed"`
						MaxBikeCadence int     `xml:"ns3:MaxBikeCadence"`
						AvgWatts       int     `xml:"ns3:AvgWatts"`
						MaxWatts       int     `xml:"ns3:MaxWatts"`
					} `xml:"ns3:LX"`
				} `xml:"Extensions"`
			} `xml:"Lap"`
		} `xml:"Activity"`
	} `xml:"Activities"`
}

type Trackpoint struct {
	Text         string `xml:",chardata"`
	Time         string `xml:"Time"`
	HeartRateBpm struct {
		Text  string `xml:",chardata"`
		Value int    `xml:"Value"`
	} `xml:"HeartRateBpm"`
	Cadence    int `xml:"Cadence"`
	Extensions struct {
		Text string `xml:",chardata"`
		TPX  struct {
			Text  string  `xml:",chardata"`
			Speed float64 `xml:"ns3:Speed"`
			Watts int     `xml:"ns3:Watts"`
		} `xml:"ns3:TPX"`
	} `xml:"Extensions"`
}
