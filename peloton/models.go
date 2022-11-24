package peloton

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Ftp            int       `json:"cycling_workout_ftp"`
	CustomMaxHr    int       `json:"customized_max_heart_rate"`
	TotalWorkOuts  int       `json:"total_workouts"`
	DefaultMaxHr   int       `json:"default_max_heart_rate"`
	EstimatedFtp   int       `json:"estimated_cycling_ftp"`
	DefaultHrZones []float64 `json:"default_heart_rate_zones"`
	Birthday       int       `json:"birthday"`
	Location       string    `json:"location"`
	Height         float64   `json:"height"`
	Weight         float64   `json:"weight"`
	Gender         string    `json:"gender"`
	LastWorkout    time.Time `json:"last_workout_at"`
}

type Workouts struct {
	Count     int `json:"count"`
	Limit     int `json:"limit,omitempty"`
	Page      int `json:"page,omitempty"`
	Total     int `json:"total,omitempty"`
	PageCount int `json:"page_count,omitempty"`
	Data      []WorkoutData
}

type Peloton struct {
	Ride Ride `json:"ride"`
}

type Ride struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Title       string `json:"title"`
}

type WorkoutData struct {
	FitnessDiscipline string  `json:"fitness_discipline"`
	Created           int     `json:"created_at"`
	EndTime           int     `json:"end_time"`
	ID                string  `json:"id"`
	StartTime         int     `json:"start_time"`
	UserID            string  `json:"user_id"`
	Status            string  `json:"status"`
	PersonalRecord    bool    `json:"is_total_work_personal_record"`
	PedalingMetrics   bool    `json:"has_pedaling_metrics"`
	Peloton           Peloton `json:"peloton"`
}

type WorkoutDetailSegmentList struct {
	ID              string  `json:"id"`
	Length          int     `json:"length"`
	StartTimeOffset int     `json:"start_time_offset"`
	IconURL         string  `json:"icon_url"`
	IntensityInMets float64 `json:"intensity_in_mets"`
	MetricsType     string  `json:"metrics_type"`
	IconName        string  `json:"icon_name"`
	IconSlug        string  `json:"icon_slug"`
	Name            string  `json:"name"`
	IsDrill         bool    `json:"is_drill"`
}

type WorkoutDetailAverageSummaries struct {
	DisplayName string  `json:"display_name"`
	DisplayUnit string  `json:"display_unit"`
	Value       float64 `json:"value"`
	Slug        string  `json:"slug"`
}

type WorkoutDetailSummaries struct {
	DisplayName string  `json:"display_name"`
	DisplayUnit string  `json:"display_unit"`
	Value       float64 `json:"value"`
	Slug        string  `json:"slug"`
}

type WorkoutDetailMetrics struct {
	DisplayName         string                     `json:"display_name"`
	DisplayUnit         string                     `json:"display_unit"`
	MaxValue            float64                    `json:"max_value"`
	AverageValue        float64                    `json:"average_value"`
	Values              []float64                  `json:"values"`
	Slug                string                     `json:"slug"`
	Zones               []WorkoutDetailMetricZones `json:"zones,omitempty"`
	MissingDataDuration int                        `json:"missing_data_duration,omitempty"`
}

type WorkoutDetailMetricZones struct {
	DisplayName string `json:"display_name"`
	Slug        string `json:"slug"`
	Range       string `json:"range"`
	Duration    int    `json:"duration"`
	MaxValue    int    `json:"max_value"`
	MinValue    int    `json:"min_value"`
}

type WorkoutDetailSplits struct {
	DistanceMarker  float64     `json:"distance_marker"`
	Order           int         `json:"order"`
	Seconds         float64     `json:"seconds"`
	ElevationChange interface{} `json:"elevation_change"`
	HasFloorSegment bool        `json:"has_floor_segment"`
	IsBest          bool        `json:"is_best"`
}

type WorkoutDetailSplitsData struct {
	DistanceMarkerDisplayUnit  string                `json:"distance_marker_display_unit"`
	ElevationChangeDisplayUnit string                `json:"elevation_change_display_unit"`
	Splits                     []WorkoutDetailSplits `json:"splits"`
}

type MetricData struct {
	Slug  string  `json:"slug"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type WorkoutSplitMetrics struct {
	IsBest          bool         `json:"is_best"`
	HasFloorSegment bool         `json:"has_floor_segment"`
	Data            []MetricData `json:"data"`
}

type WorkoutSplitMetricHeader struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"display_name"`
}

type WorkoutSplitsMetrics struct {
	Header  []WorkoutSplitMetricHeader `json:"header"`
	Metrics []WorkoutSplitMetrics      `json:"metrics"`
}

type WorkdayTargetMetricsPerfData struct {
	TargetMetrics []interface{} `json:"target_metrics"`
	TimeInMetric  []interface{} `json:"time_in_metric"`
}

type EffortZoneHeartRateDurations struct {
	HeartRateZ1Duration int `json:"heart_rate_z1_duration"`
	HeartRateZ2Duration int `json:"heart_rate_z2_duration"`
	HeartRateZ3Duration int `json:"heart_rate_z3_duration"`
	HeartRateZ4Duration int `json:"heart_rate_z4_duration"`
	HeartRateZ5Duration int `json:"heart_rate_z5_duration"`
}

type EffortZones struct {
	TotalEffortPoints      float64                      `json:"total_effort_points"`
	HeartRateZoneDurations EffortZoneHeartRateDurations `json:"heart_rate_zone_durations"`
}

type WorkoutDetail struct {
	Title                        string
	Description                  string
	ID                           string
	FitnessDiscipline            string
	DataGranularityInSeconds     int
	StartTime                    time.Time
	EndTime                      time.Time
	Duration                     int                             `json:"duration"`
	IsClassPlanShown             bool                            `json:"is_class_plan_shown"`
	SegmentList                  []WorkoutDetailSegmentList      `json:"segment_list"`
	SecondsSincePedalingStart    []int                           `json:"seconds_since_pedaling_start"`
	AverageSummaries             []WorkoutDetailAverageSummaries `json:"average_summaries"`
	Summaries                    []WorkoutDetailSummaries        `json:"summaries"`
	Metrics                      []WorkoutDetailMetrics          `json:"metrics"`
	SplitsData                   WorkoutDetailSplitsData         `json:"splits_data"`
	SplitsMetrics                WorkoutSplitsMetrics            `json:"splits_metrics"`
	TargetMetricsPerformanceData WorkdayTargetMetricsPerfData    `json:"target_metrics_performance_data"`
	EffortZones                  EffortZones                     `json:"effort_zones"`
}
