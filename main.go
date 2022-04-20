package main

import (
	"dyngo/config"
	"dyngo/detection"
	"dyngo/logger"
	"dyngo/services"
	"fmt"

	"github.com/robfig/cron/v3"
)

var currentIPv4 string
var currentIPv6 string

func main() {
	c := cron.New(cron.WithSeconds())

	fmt.Println("")
	fmt.Println("===========================")
	fmt.Println("==   Welcome to DynGO!   ==")
	fmt.Println("==   Version: 0.0.1      ==")
	fmt.Println("===========================")
	fmt.Println("")

	config.Parse()

	logger.LogLevel = config.LogLevel
	logger.Info.Printf("Initiating cron job with pattern %v\n", config.Cron)

	// I should probably loop over Services, but it's a struct and I don't know
	// what golangs equivalent to Object.keys() is...
	if config.Services.Desec.Domains != nil {
		services.Register(services.NewDesec(config.Services.Desec))
	}

	c.AddFunc(config.Cron, updateDynDNS)

	updateDynDNS() // Run immediatly once
	c.Run()        // Run cron
}

func updateDynDNS() {
	var upstreamIPv4 string
	var upstreamIPv6 string

	if services.AtLeastOneDomainRequires("v4") {
		upstreamIPv4 = detection.GetIPv4()

		if upstreamIPv4 == "" {
			logger.Error.Printf("Could not determine IPv4, skipping...")
			return
		}

		if currentIPv4 != upstreamIPv4 {
			logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
			currentIPv4 = upstreamIPv4
		}

	}

	if services.AtLeastOneDomainRequires("v6") {
		upstreamIPv6 = detection.GetIPv6()

		if upstreamIPv6 == "" {
			logger.Error.Printf("Could not determine IPv6, skipping...")
			return
		}

		if currentIPv6 != upstreamIPv6 {
			logger.Info.Printf("Detected change in IPv6 Address: '%v' -> '%v' \n", currentIPv6, upstreamIPv6)
			currentIPv6 = upstreamIPv6
		}
	}

	for _, service := range services.Registered {
		service.SetTargetIPv4(upstreamIPv4)
		service.SetTargetIPv6(upstreamIPv6)
		service.UpdateAllDomains()
	}
}
