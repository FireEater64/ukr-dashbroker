package dashbroker

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	log "github.com/cihub/seelog"
)

type Config struct {
	DatabaseType             string `yaml:"databaseType"`
	DatabaseConnectionString string `yaml:"databaseConnectionString"`
	ClockworkSMSApiKey       string `yaml:"clockworkApiKey"`
}

var Configuration Config

func LoadConfiguration(fileName string) {
	configFile, fileErr := ioutil.ReadFile(fileName)

	if fileErr != nil {
		log.Critical("Error whilst reading config file at %s: %s", fileName, fileErr)
		panic(fileErr)
	}

	yamlErr := yaml.Unmarshal(configFile, &Configuration)

	if yamlErr != nil {
		log.Critical("Error whilst parsing yaml: %s", yamlErr)
		panic(yamlErr)
	}
}
