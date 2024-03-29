package main

import (
	"dyngo/config"
	"dyngo/detection/triggers"
	"dyngo/logger"
	"dyngo/services"
	"flag"
	"fmt"
	"os"
	"strings"
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

	if len(config.Services.Generic) > 0 {
		for _, service := range config.Services.Generic {
			services.Register(services.NewGenericService(*service))
		}
	}

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

	triggers.SetupTriggers()

	logger.Info.Println("No continuous triggers configured. Shutting down, goodbye!")
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
