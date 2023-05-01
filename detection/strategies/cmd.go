package strategies

import (
	"os/exec"
	"strings"
)

type CmdDetectionStrategy struct {
	BaseDetectionStrategy
	Command string
}

func NewCmdDetectionStrategy(Command string) DetectionStrategy {
	return &CmdDetectionStrategy{
		BaseDetectionStrategy: NewBaseDetectionStrategy("cmd"),
		Command:               Command,
	}
}

func (strategy *CmdDetectionStrategy) Execute() string {
	strategy.Logger.Debug.Printf("Executing command: " + strategy.Command)

	cmd := exec.Command("/bin/sh", "-c", strategy.Command)

	out, err := cmd.Output()

	if err != nil {
		strategy.Logger.Error.Printf("Detection failed: %s", err)

		return ""
	}

	ip := strings.TrimSpace(string(out))

	strategy.Logger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
