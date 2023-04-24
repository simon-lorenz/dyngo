package main

import (
	"dyngo/detection"
	"dyngo/logger"
	"dyngo/services"
)

var currentIPv4 string = "Unkown"
var currentIPv6 string = "Unkown"

func runDynDNSUpdater() {
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
