package database

var configuration Configuration

type Configuration struct {
	DatabaseType             string
	DatabaseConnectionString string
}

func SetConfiguration(givenConfiguration Configuration) {
	configuration = givenConfiguration
}
