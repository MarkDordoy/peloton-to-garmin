package version

var (
	version = "develop"
)

func GetVersion() string {
	return version
}

func GetUA() string {
	return "peloton-to-garmin/" + version
}
