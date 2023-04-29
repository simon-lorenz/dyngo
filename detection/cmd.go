package detection

import (
	"dyngo/logger"
	"os/exec"
	"strings"
)

var CmdDetectionLogger = logger.NewLoggerCollection("detection/cmd")

func getIpAddressFromCmd(s string) string {
	CmdDetectionLogger.Debug.Printf("Executing command: " + s)

	cmd := exec.Command("/bin/sh", "-c", s)

	out, err := cmd.Output()

	if err != nil {
		CmdDetectionLogger.Error.Printf("Detection failed: %s", err)

		return ""
	}

	ip := strings.TrimSpace(string(out))

	CmdDetectionLogger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
