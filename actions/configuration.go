package actions

var configuration Configuration

type Configuration struct {
	ClockworkAPIKey string
	YoAPIKey        string
}

func SetConfiguration(givenConfiguration Configuration) {
	configuration = givenConfiguration
}
