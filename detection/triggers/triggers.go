package triggers

import (
	"dyngo/config"
	"dyngo/detection"
	"dyngo/logger"
	"sync"

	"github.com/robfig/cron/v3"
)

var triggersLogger = logger.NewLoggerCollection("detection/triggers")

func SetupTriggers() {
	var wg sync.WaitGroup

	if config.Detection.Triggers.Cron != "" {
		triggersLogger.Info.Printf("Initiating cron trigger with pattern %v\n", config.Detection.Triggers.Cron)
		wg.Add(1)
		go runCronTrigger(config.Detection.Triggers.Cron)
	}

	if config.Detection.Triggers.Startup {
		triggersLogger.Info.Printf("Initiating startup trigger")
		runStartupTrigger()
	}

	wg.Wait()
}

func runCronTrigger(pattern string) {
	cron := cron.New(cron.WithSeconds())
	cron.AddFunc(pattern, func() { detection.RunDetection("cron") })
	cron.Run()
}

func runStartupTrigger() {
	detection.RunDetection("startup")
}
