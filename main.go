package main

import (
	"dyngo/config"
	"dyngo/logger"
	"dyngo/services"
	"flag"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
)

type Flags struct {
	config *string
}

var version = "development"

func main() {
	var flags Flags = setupAndParseFlags()

	printWelcomeMessage()

	config.Parse(*flags.config)
	logger.SetLogLevel(config.Log.Level)
	logger.Info.Println("Using configuration file " + *flags.config)

	// I should probably loop over Services, but it's a struct and I don't know
	// what golangs equivalent to Object.keys() is...
	if config.Services.Desec.Domains != nil {
		services.Register(services.NewDesec())
	}

	// Run cron
	logger.Info.Printf("Initiating cron job with pattern %v\n", config.Cron)
	c := cron.New(cron.WithSeconds())
	c.AddFunc(config.Cron, runDynDNSUpdater)
	defer c.Run()

	// Run once immediatly
	runDynDNSUpdater()
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
	fmt.Printf("==   Version: %-15s   ==\n", version)
	fmt.Printf("%s\n", strings.Repeat("=", 34))
}
