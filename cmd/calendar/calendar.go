package main

import (
	"flag"
	"log"
	"time"

	"github.com/blang/whoshome"
	"github.com/blang/whoshome/calendar"
)

var (
	flagClientSecretFile = flag.String("client-secret-file", "client-secret.json", "OAuth2 token file")
	flagTokenFile        = flag.String("client-token-file", "client-token.json", "Google OAuth2 token file, generate: -init-token")
	flagARPFile          = flag.String("arp-file", "/proc/net/arp", "ARP linux file")
	flagInitToken        = flag.Bool("init-token", false, "Create token file")
	flagTickTime         = flag.Duration("interval", time.Minute, "Update interval")
	flagCalendarID       = flag.String("calendar-id", "", "Google calendar id e.g. abc@group.calendar.google.com")
	flagMAC              = flag.String("mac", "00:01:02:aa:bb:cc", "MAC to identify presence")
	flagEventTitle       = flag.String("event-title", "Present", "Event title")
	flagEventColorID     = flag.String("event-color-id", "11", "Event colorid")
	flagGraceThreshold   = flag.Int("grace-threshold", 5, "Threshold how many ticks no present triggers an event")
)

func main() {
	flag.Parse()

	arp := whoshome.NewARPProvider(*flagARPFile, map[string]string{*flagMAC: *flagEventTitle})
	cal := &calendar.Calendar{
		Provider:         arp,
		TickTime:         *flagTickTime,
		CalendarID:       *flagCalendarID,
		EventColorID:     *flagEventColorID,
		ClientSecretFile: *flagClientSecretFile,
		ClientTokenFile:  *flagTokenFile,
		EventTitle:       *flagEventTitle,
		Threshold:        *flagGraceThreshold,
		PresentFunc: func(present []string) bool {
			for _, v := range present {
				if v == *flagEventTitle {
					return true
				}
			}
			return false
		},
	}
	var err error
	if *flagInitToken {
		err = cal.InitAuth()
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		log.Printf("InitAuth successful")
	}
	err = cal.Init()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	log.Printf("Init successful")

	cal.Run()
}
