package actions

import (
	"sync"

	"github.com/FireEater64/clockwork-go"
	log "github.com/cihub/seelog"
)

var clock *clockwork.Clockwork

func SendSMS(recipient string, message string) {
	checkClockwork()
	toSend := clockwork.SMS{To: recipient, Message: message}
	messageResponse := clock.SendSMS(toSend)
	if messageResponse.SMSResult[0].ErrorMessage != "" {
		log.Warnf("Error sending SMS to %s: %s", recipient, messageResponse.SMSResult[0].ErrorMessage)
	}
}

func SendSMSMultiple(recipients []string, message string) {
	for _, recipient := range recipients {
		SendSMS(recipient, message)
	}
}

func SendSMSAsync(recipient string, message string) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()
	go SendSMS(recipient, message)
	return &wg
}

func SendSMSMultipleAsync(recipients []string, message string) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(len(recipients))

	for _, recipient := range recipients {
		go func(recipient string, message string) {
			defer wg.Done()
			SendSMS(recipient, message)
		}(recipient, message)
	}

	return &wg
}

func checkClockwork() {
	if clock == nil {
		clock = clockwork.NewClockwork(configuration.ClockworkAPIKey)
	}
}
