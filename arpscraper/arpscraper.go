package arpscraper

import (
	"net"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/mdlayher/arp"
)

type ARPScraper struct {
	macChan chan string
	wg      sync.WaitGroup
	run     bool
}

func NewScraper(MACChan chan string) *ARPScraper {
	toReturn := ARPScraper{}
	toReturn.macChan = MACChan
	toReturn.wg = sync.WaitGroup{}
	return &toReturn
}

func (as *ARPScraper) Start() {
	log.Debug("Start scraping ARP messages")
	as.wg.Add(1)
	as.run = true
	go as.listenForArpMessages()
}

func (as *ARPScraper) Stop() {
	log.Debug("ARP Scraper stopping")
	as.run = false
	as.wg.Wait()
}

func (as *ARPScraper) listenForArpMessages() {
	defer as.wg.Done()

	iface, ifaceErr := net.InterfaceByName("eth0")
	if ifaceErr != nil {
		log.Criticalf("Error obtaining interface: %s", ifaceErr.Error())
		panic(ifaceErr)
	}

	arpClient, clientErr := arp.NewClient(iface)
	if clientErr != nil {
		log.Criticalf("Error obtaining interface: %s", clientErr.Error())
		panic(clientErr)
	}

	for as.run {
		p, _, err := arpClient.Read()
		if err != nil {
			log.Criticalf("Read error: %s", err)
			time.Sleep(3 * time.Second)
			continue
		}

		if p.Operation == arp.OperationRequest {
			as.macChan <- p.SenderHardwareAddr.String()
		}
	}
}
