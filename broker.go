package main

import (
	"flag"
	"io/ioutil"

	"github.com/FireEater64/ukrdashbroker/actions"
	"github.com/FireEater64/ukrdashbroker/arpscraper"
	"github.com/FireEater64/ukrdashbroker/database"
	log "github.com/cihub/seelog"
	"gopkg.in/yaml.v2"
)

var configFile string

func main() {
	initializeLogging()
	loadFlags()
	loadConfig()

	defer log.Flush()

	macAddressChannel := make(chan string, 10)

	scraper := arpscraper.NewScraper(macAddressChannel)
	scraper.Start()
	for {
		macAddress := <-macAddressChannel
		go processMacAddress(macAddress)
	}
}

func processMacAddress(givenMacAddress string) {
	switch givenMacAddress {
	case "74:75:48:2e:2b:4c": // Dinner is ready
		log.Debug("SmartWater Pressed")
		housematesNumbers := database.GetAllActiveHousematesNumbers()
		wg := actions.SendSMSMultipleAsync(housematesNumbers, "Dinner is ready!")
		go actions.YoFromHome() // We don't really care if Yo suceeds or not
		database.LogButtonPress(givenMacAddress, "DinnerNotification")
		wg.Wait()
	case "74:c2:46:84:ab:8e": // Nyan cat time
		log.Debug("Gilette Pressed")
		actions.PlayNyanCatInDiningRoom()
		database.LogButtonPress(givenMacAddress, "NyanCatTime")
	}
}

func loadFlags() {
	flag.StringVar(&configFile, "config", "config.yml", "The path to the config.yml file")
	flag.Parse()
}

func loadConfig() {
	mainConfig := loadMasterConfigurationFromFile(configFile)

	// Init database
	dbConfig := database.Configuration{}
	dbConfig.DatabaseType = mainConfig.DatabaseType
	dbConfig.DatabaseConnectionString = mainConfig.DatabaseConnectionString
	database.SetConfiguration(dbConfig)

	// Init actions
	actionsConfig := actions.Configuration{}
	actionsConfig.ClockworkAPIKey = mainConfig.ClockworkSMSAPIKey
	actionsConfig.YoAPIKey = mainConfig.YoAPIKey
	actionsConfig.YoLocation = mainConfig.YoLocation
	actions.SetConfiguration(actionsConfig)
}

func initializeLogging() {
	logger, err := log.LoggerFromConfigAsFile("logconfig.xml")

	if err != nil {
		log.Criticalf("An error occurred whilst initializing logging\n", err.Error())
		panic(err)
	}

	log.ReplaceLogger(logger)
}

type Config struct {
	DatabaseType             string `yaml:"databaseType"`
	DatabaseConnectionString string `yaml:"databaseConnectionString"`
	ClockworkSMSAPIKey       string `yaml:"clockworkApiKey"`
	YoAPIKey                 string `yaml:"yoApiKey"`
	YoLocation               string `yaml:"yoLocation"`
}

func loadMasterConfigurationFromFile(fileName string) *Config {
	configFile, fileErr := ioutil.ReadFile(fileName)

	if fileErr != nil {
		log.Critical("Error whilst reading config file at %s: %s", fileName, fileErr)
		panic(fileErr)
	}

	toReturn := Config{}

	yamlErr := yaml.Unmarshal(configFile, &toReturn)

	if yamlErr != nil {
		log.Critical("Error whilst parsing yaml: %s", yamlErr)
		panic(yamlErr)
	}

	return &toReturn
}
