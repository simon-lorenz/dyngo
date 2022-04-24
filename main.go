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

type DynDnsService interface {
	UpdateIPv4(string)
	UpdateIPv6(string)
	GetHosts() []config.HostConfiguration
	GetName() string
}

var activeServices []DynDnsService

func main() {
	c := cron.New(cron.WithSeconds())

	fmt.Println("")
	fmt.Println("===========================")
	fmt.Println("==   Welcome to DynGO!   ==")
	fmt.Println("==   Version: 0.0.1      ==")
	fmt.Println("===========================")
	fmt.Println("")

	config.Parse()

	// I should probably loop over Services, but it's a struct and I don't know
	// what golangs equivalent to Object.keys() is...
	if config.Services.Desec.Hosts != nil {
		registerService(clients.NewDesec(config.Services.Desec))
	}

	logger.Info.Printf("Initiating cron job with pattern %v\n", config.Cron)

	c.AddFunc(config.Cron, updateDynDNS)

	updateDynDNS() // Run immediatly

	c.Run() // Run cron
}

func atLeastOneHostRequests(protocol string) bool {
	for _, service := range activeServices {
		for _, host := range service.GetHosts() {
			if (protocol == "v4" && host.V4) || (protocol == "v6" && host.V6) {
				return true
			}
		}
	}

	return false
}

func registerService(service DynDnsService) {
	activeServices = append(activeServices, service)
	logger.Info.Printf("Registered service '%v'", service.GetName())
}

func updateDynDNS() {
	var upstreamIPv4 string
	var upstreamIPv6 string

	if atLeastOneHostRequests("v4") {
		upstreamIPv4 = detection.GetIPv4()
	}

	if atLeastOneHostRequests("v6") {
		upstreamIPv6 = detection.GetIPv6()
	}

	for _, service := range activeServices {
		for _, host := range service.GetHosts() {
			if host.V4 && currentIPv4 != upstreamIPv4 {
				logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
				service.UpdateIPv4(upstreamIPv4)
				currentIPv4 = upstreamIPv4
			}

			if host.V6 && currentIPv6 != upstreamIPv6 {
				logger.Info.Printf("Detected change in IPv6 Address: '%v' -> '%v' \n", currentIPv6, upstreamIPv6)
				service.UpdateIPv6(upstreamIPv6)
				currentIPv6 = upstreamIPv6
			}
		}
	}
}
