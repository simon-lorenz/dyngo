package strategies

import (
	"dyngo/logger"
	"os/exec"
	"strings"
)

var cmdDetectionLogger = logger.NewLoggerCollection("detection/strategies/cmd")

func GetIpAddressFromCmd(s string) string {
	cmdDetectionLogger.Debug.Printf("Executing command: " + s)

	cmd := exec.Command("/bin/sh", "-c", s)

	out, err := cmd.Output()

	if err != nil {
		cmdDetectionLogger.Error.Printf("Detection failed: %s", err)

		return ""
	}

	ip := strings.TrimSpace(string(out))

	cmdDetectionLogger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
