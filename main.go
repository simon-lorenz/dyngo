package main

import (
	"dyngo/clients"
	"dyngo/config"
	"dyngo/detection"
	"dyngo/logger"
	"fmt"

	"github.com/robfig/cron/v3"
)

var currentIPv4 string
var currentIPv6 string

func main() {
	c := cron.New(cron.WithSeconds())

	fmt.Println("===========================")
	fmt.Println("==   Welcome to DynGO!   ==")
	fmt.Println("==   Version: 0.0.1      ==")
	fmt.Println("===========================")

	fmt.Println("")

	config.Parse()

	logger.Info.Printf("Initiating cron job with pattern %v\n", config.Cron)

	c.AddFunc(config.Cron, updateDynDNS)

	updateDynDNS() // Run immediatly

	c.Run() // Run cron
}

func updateDynDNS() {
	upstreamIPv4 := detection.GetIPv4()
	// upstreamIPv6 := detection.GetIPv6()
	upstreamIPv6 := ""

	// TODO: Register active services of configuration file and loop over them

	if config.Services.Desec.Username != "" { // TODO: Normally I would check for nil but this doesn't work
		for _, host := range config.Services.Desec.Hosts {
			desec := clients.NewDesec(config.Services.Desec.Username, config.Services.Desec.Password)

			if host.V4 && currentIPv4 != upstreamIPv4 {
				logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
				desec.UpdateIPv4(detection.GetIPv4(), host.Host)
				currentIPv4 = upstreamIPv4
			}

			if host.V6 && currentIPv6 != upstreamIPv6 {
				logger.Info.Printf("Detected change in IP6 Address: '%v' -> '%v' \n", currentIPv6, upstreamIPv6)
				desec.UpdateIPv6(detection.GetIPv6(), host.Host)
				currentIPv6 = upstreamIPv6
			}
		}
	}
}
