package main

import (
	"dyngo/clients"
	"dyngo/config"
	"dyngo/helper"
	"dyngo/logger"
	"fmt"
	"os"

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

	if !config.IPv4Enabled && !config.IPv6Enabled {
		logger.Error.Println("Neither IPv4 nor IPv6 updates are enabled.")
		os.Exit(1)
	}

	c.AddFunc(config.Cron, updateDynDNS)

	updateDynDNS() // Run immediatly

	c.Run() // Run cron
}

func updateDynDNS() {

	services := []clients.DynDnsService{clients.NewDesec()}

	if config.IPv4Enabled {
		var upstreamIPv4 = helper.GetIPv4()

		if upstreamIPv4 != currentIPv4 {
			logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
			services[0].UpdateIPv4(upstreamIPv4)
			currentIPv4 = upstreamIPv4
		}
	}

	if config.IPv6Enabled {
		var upstreamIPv6 = helper.GetIPv6()

		if upstreamIPv6 != currentIPv6 {
			logger.Info.Printf("Detected change in IP6 Address: '%v' -> '%v' \n", currentIPv6, upstreamIPv6)
			services[0].UpdateIPv6(upstreamIPv6)
			currentIPv6 = upstreamIPv6
		}
	}
}
