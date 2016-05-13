package actions

import (
	"github.com/cemkiy/go-yo" // Used until location feature is merged into master
)

var client *yo.Client

func YoAll() {
	checkClient()

	client.YoAll()
}

func YoFromHome() {
	checkClient()

	client.YoAllLocation(configuration.YoLocation)
}

func checkClient() {
	if client == nil {
		client = yo.NewClient(configuration.YoAPIKey)
	}
}
