package main

import (
	"dyngo/config"
	"dyngo/logger"
	"dyngo/services"
	"flag"
	"fmt"
	"os"
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

	if config.Services.Desec != nil {
		services.Register(services.NewDesec())
	}

	if config.Services.Porkbun != nil {
		services.Register(services.NewPorkbun())
	}

	if len(services.Registered) == 0 {
		logger.Fatal.Println("No services registered")
		os.Exit(1)
	}

	// Run cron
	if config.Detection.Triggers.Cron != "" {
		logger.Info.Printf("Initiating cron job with pattern %v\n", config.Detection.Triggers.Cron)
		c := cron.New(cron.WithSeconds())
		c.AddFunc(config.Detection.Triggers.Cron, runDynDNSUpdater)
		defer c.Run()
	}

	// Run once immediatly
	if config.Detection.Triggers.Startup {
		runDynDNSUpdater()
	}

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
