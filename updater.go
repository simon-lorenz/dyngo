package main

import (
	"dyngo/detection"
	"dyngo/logger"
	"dyngo/services"
)

var currentIPv4 string = "Unknown"
var currentIPv6 string = "Unknown"

func runDynDNSUpdater() {
	var upstreamIPv4 string
	var upstreamIPv6 string
	var err error

	if services.AtLeastOneDomainRequires("v4") {
		upstreamIPv4, err = detection.GetIPv4()

		if err != nil {
			logger.Error.Println(err.Error())
			logger.Error.Println("Skipping update")
			return
		}

		if currentIPv4 != upstreamIPv4 {
			logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
			currentIPv4 = upstreamIPv4
		}
	}

	if services.AtLeastOneDomainRequires("v6") {
		upstreamIPv6, err = detection.GetIPv6()

		if err != nil {
			logger.Error.Println(err.Error())
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
