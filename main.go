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

	fmt.Println("")
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
	IPv4Enabled := config.AtLeastOneIPv4UpdateRequested()
	IPv6Enabled := config.AtLeastOneIPv6UpdateRequested()

	if IPv4Enabled {
		upstreamIPv4 := detection.GetIPv4()

		if config.Services.Desec.Hosts != nil {
			for _, host := range config.Services.Desec.Hosts {
				desec := clients.NewDesec(config.Services.Desec.Username, config.Services.Desec.Password)

				if host.V4 && currentIPv4 != upstreamIPv4 {
					logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", currentIPv4, upstreamIPv4)
					desec.UpdateIPv4(detection.GetIPv4(), host.Host)
					currentIPv4 = upstreamIPv4
				}

			}
		}
	}

	if IPv6Enabled {
		upstreamIPv6 := detection.GetIPv6()

		// TODO: Register active services of configuration file and loop over them
		if config.Services.Desec.Hosts != nil {
			for _, host := range config.Services.Desec.Hosts {
				desec := clients.NewDesec(config.Services.Desec.Username, config.Services.Desec.Password)

				if host.V6 && currentIPv6 != upstreamIPv6 {
					logger.Info.Printf("Detected change in IP6 Address: '%v' -> '%v' \n", currentIPv6, upstreamIPv6)
					desec.UpdateIPv6(detection.GetIPv6(), host.Host)
					currentIPv6 = upstreamIPv6
				}
			}
		}
	}
}
