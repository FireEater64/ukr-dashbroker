package dashbroker

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "github.com/cihub/seelog"
)

type Config struct {
	DatabaseType             string `yaml:"databaseType"`
	DatabaseConnectionString string `yaml:"databaseConnectionString"`
	ClockworkSMSApiKey       string `yaml:"clockworkApiKey"`
}

var configuration Config

func LoadConfiguration(fileName string) {
	configFile, fileErr := ioutil.ReadFile(fileName)

	if fileErr != nil {
		log.Critical("Error whilst reading config file at %s: %s", fileName, fileErr)
		panic(fileErr)
	}

	yamlErr := yaml.Unmarshal(configFile, &configuration)

	if yamlErr != nil {
		log.Critical("Error whilst parsing yaml: %s", yamlErr)
		panic(yamlErr)
	}
}
