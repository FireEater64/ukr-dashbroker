package actions

import (
	"github.com/cemkiy/go-yo" // Used until location feature is merged into master
)

const ukrCoordinates string = "53.4551468,-2.2108870"

var client *yo.Client

func YoAll() {
	checkClient()

	client.YoAll()
}

func YoFromHome() {
	checkClient()

	client.YoAllLocation(ukrCoordinates)
}

func checkClient() {
	if client == nil {
		client = yo.NewClient(configuration.YoAPIKey)
	}
}
