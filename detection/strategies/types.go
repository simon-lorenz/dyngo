package strategies

import "dyngo/logger"

type DetectionStrategy interface {
	Execute() string
}

type BaseDetectionStrategy struct {
	Name   string
	Logger *logger.LoggerCollection
}

func NewBaseDetectionStrategy(name string) BaseDetectionStrategy {
	return BaseDetectionStrategy{
		Name:   name,
		Logger: logger.NewLoggerCollection("detection/strategies/" + name),
	}
}
