package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/FireEater64/clockwork-go"
	"github.com/FireEater64/ukr-dashbroker"
	log "github.com/cihub/seelog"
	"github.com/mdlayher/arp"
)

var wg sync.WaitGroup
var clock *clockwork.Clockwork

func main() {
	initializeLogging()

	defer log.Flush()

	dashbroker.LoadConfiguration("config.yml")
	clock = clockwork.NewClockwork(dashbroker.Configuration.ClockworkSMSApiKey)

	macAddressChannel := make(chan string, 10)

	wg.Add(2)
	go listenForButtonPress(macAddressChannel)
	go listenForMacAddresses(macAddressChannel)
	wg.Wait()
}

func initializeLogging() {
	logger, err := log.LoggerFromConfigAsFile("logconfig.xml")

	if err != nil {
		log.Criticalf("An error occurred whilst initializing logging\n", err.Error())
		panic(err)
	}

	log.ReplaceLogger(logger)
}

func listenForMacAddresses(inChan chan string) {
	buttonAddresses := make(map[string]bool, 3)
	buttons := dashbroker.GetAllButtons()
	for _, button := range buttons {
		buttonAddresses[button.MacAddress] = true
	}

	log.Debugf("Loaded %d buttons. Listening for button presses.", len(buttonAddresses))

	for {
		receivedAddress := <-inChan

		if buttonAddresses[receivedAddress] {
			switch receivedAddress {
			case "74:75:48:2e:2b:4c": // Kitchen button
				tellHouseDinnerIsReady()
				dashbroker.LogButtonPress(receivedAddress, "DinnerNotification")
			case "74:c2:46:84:ab:8e": // YouTube button
				playNyanCatInDiningRoom()
				dashbroker.LogButtonPress(receivedAddress, "Nyan Cat")
			default:
				dashbroker.LogButtonPress(receivedAddress, "Debug")
			}
		} else {
			log.Debugf("Unknown MAC address: %s", receivedAddress)
		}

		log.Debug(receivedAddress)
	}
}

func playNyanCatInDiningRoom() {
	url := "http://192.168.1.33:8008/apps/YouTube"
	requestBody := strings.NewReader("v=QH2-TGUlwu4")

	resp, respErr := http.Post(url, "text/plain", requestBody)
	if respErr != nil {
		log.Errorf("Error sending TV request: %s", respErr)
	}
	defer resp.Body.Close()

	body, bodyErr := ioutil.ReadAll(resp.Body)
	if bodyErr != nil {
		log.Errorf("Error whilst reading body: %s", bodyErr)
	}

	log.Debug("Response from TV: %s", string(body))
}

func tellHouseDinnerIsReady() {

	wg := sync.WaitGroup{}

	houseMates := dashbroker.GetAllActiveHousemates()
	for _, housemate := range houseMates {
		log.Debugf("Sending SMS to: %s", housemate.FirstName)
		wg.Add(1)
		go sendSMSAsync(housemate.PhoneNumber, "Dinner is ready!", &wg)
	}

	wg.Wait()
	log.Debug("Finished sending SMS messages")
}

func sendSMS(recipient string, message string) {
	toSend := clockwork.SMS{To: recipient, Message: message}
	messageResponse := clock.SendSMS(toSend)
	if messageResponse.SMSResult[0].ErrorMessage != "" {
		log.Warnf("Error sending SMS to %s: %s", recipient, messageResponse.SMSResult[0].ErrorMessage)
	}
}

func sendSMSAsync(recipient string, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	sendSMS(recipient, message)
}

func listenForButtonPress(outChannel chan string) {
	iface, ifaceErr := net.InterfaceByName("en0")
	if ifaceErr != nil {
		log.Criticalf("Error obtaining interface: %s", ifaceErr.Error())
		panic(ifaceErr)
	}

	arpClient, clientErr := arp.NewClient(iface)
	if clientErr != nil {
		log.Criticalf("Error obtaining interface: %s", clientErr.Error())
		panic(clientErr)
	}

	log.Debug("Listening")

	for {
		p, _, err := arpClient.Read()
		if err != nil {
			log.Criticalf("Read error: %s", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if p.Operation == arp.OperationRequest &&
			p.SenderIP.Equal(net.IPv4zero) {
			outChannel <- p.SenderHardwareAddr.String()
		}
	}
}
