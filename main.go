package main

import (
	"dyngo/config"
	"dyngo/detection"
	"dyngo/logger"
	"dyngo/services"
	"flag"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

var currentIPv4 string
var currentIPv6 string

type Flags struct {
	config *string
}

func main() {
	c := cron.New(cron.WithSeconds())

	var flags Flags = setupAndParseFlags()

	printWelcomeMessage()

	config.Parse(*flags.config)

	logger.SetLogLevel(config.LogLevel)
	logger.Info.Println("Using configuration file " + *flags.config)

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

func setupAndParseFlags() Flags {
	var flags Flags

	flags.config = flag.String("config", "/etc/dyngo/config.yaml", "path to configuration file")

	flag.Parse()

	return flags
}

func printWelcomeMessage() {
	fmt.Printf("%s\n", strings.Repeat("=", 34))
	fmt.Printf("==   Welcome to DynGO! %s   ==\n", strings.Repeat(" ", 6))
	fmt.Printf("==   Version: %-15s   ==\n", "0.0.0")
	fmt.Printf("%s\n", strings.Repeat("=", 34))
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
