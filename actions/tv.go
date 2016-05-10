package actions

import (
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/cihub/seelog"
)

const (
	tvURL          string = "http://192.168.1.33:8008/apps/YouTube"
	nyanCatVideoID string = "v=QH2-TGUlwu4"
)

func PlayNyanCatInDiningRoom() {
	log.Debug("It's Nyan Cat time!")
	requestBody := strings.NewReader(nyanCatVideoID)

	resp, respErr := http.Post(tvURL, "text/plain", requestBody)
	if respErr != nil {
		log.Errorf("Error sending TV request: %s", respErr)
	}
	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Errorf("Error whilst reading body: %s", bodyErr)
	}

	log.Debugf("Response from TV: %s", string(body))
}
