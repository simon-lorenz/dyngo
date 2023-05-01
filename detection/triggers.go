package detection

import (
	"dyngo/config"
	"sync"

	"github.com/robfig/cron/v3"
)

func SetupTriggers() {
	var wg sync.WaitGroup

	if config.Detection.Triggers.Cron != "" {
		DetectionLogger.Info.Printf("Initiating cron trigger with pattern %v\n", config.Detection.Triggers.Cron)
		wg.Add(1)
		go RunCronTrigger(config.Detection.Triggers.Cron)
	}

	if config.Detection.Triggers.Startup {
		DetectionLogger.Info.Printf("Initiating startup trigger")
		RunStartupTrigger()
	}

	wg.Wait()
}

func RunCronTrigger(pattern string) {
	cron := cron.New(cron.WithSeconds())
	cron.AddFunc(pattern, RunDetection)
	cron.Run()
}

func RunStartupTrigger() {
	RunDetection()
}
