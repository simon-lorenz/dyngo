package main

import (
	"dyngo/detection"
	"dyngo/logger"
	"dyngo/services"
)

func runDynDNSUpdater() {
	var IPv4Changed bool = false
	var IPv6Changed bool = false
	var err error

	if services.AtLeastOneDomainRequires("v4") {
		IPv4Changed, err = detection.RefreshIPv4()

		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
	}

	if services.AtLeastOneDomainRequires("v6") {
		IPv6Changed, err = detection.RefreshIPv6()

		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
	}

	for _, service := range services.Registered {
		if IPv4Changed || IPv6Changed {
			// TODO: Split into UpdateIPv4 and UpdateIPv6
			service.UpdateAllDomains(detection.CurrentIPv4, detection.CurrentIPv6)
		}

	}
}
