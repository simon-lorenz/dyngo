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
			IPv4Changed = false
		}
	}

	if services.AtLeastOneDomainRequires("v6") {
		IPv6Changed, err = detection.RefreshIPv6()

		if err != nil {
			logger.Error.Println(err.Error())
			IPv6Changed = false
		}
	}

	for _, service := range services.Registered {
		if IPv4Changed {
			service.UpdateIPv4(detection.CurrentIPv4)
		}

		if IPv6Changed {
			service.UpdateIPv6(detection.CurrentIPv6)
		}
	}
}
