package detection

import (
	"dyngo/config"
	"dyngo/logger"
	"errors"
)

var DetectionLogger = logger.NewLoggerCollection("detection")
var CurrentIPv4 = ""
var CurrentIPv6 = ""

func RefreshIPv4() (bool, error) {
	DetectionLogger.Debug.Println("Running IPv4 detection")

	var UpstreamIPv4 string = ""

	if config.Detection.V4.Web != "" {
		UpstreamIPv4 = getIpAddressFromExternalService(config.Detection.V4.Web)
	} else if config.Detection.V4.Cmd != "" {
		UpstreamIPv4 = getIpAddressFromCmd(config.Detection.V4.Cmd)
	} else {
		return false, errors.New("Cannot determine IPv4 because no detection mechanisms are configured")
	}

	if CurrentIPv4 != UpstreamIPv4 {
		DetectionLogger.Info.Printf("Detected change in IPv4 Address: '%s' -> '%s' ", CurrentIPv4, UpstreamIPv4)
		CurrentIPv4 = UpstreamIPv4
		return true, nil
	} else {
		DetectionLogger.Debug.Println("No IPv4 change detected")
		return false, nil
	}

}

func RefreshIPv6() (bool, error) {
	DetectionLogger.Debug.Println("Running IPv6 detection")

	var UpstreamIPv6 string = ""

	if config.Detection.V6.Web != "" {
		UpstreamIPv6 = getIpAddressFromExternalService(config.Detection.V6.Web)
	} else if config.Detection.V6.Cmd != "" {
		UpstreamIPv6 = getIpAddressFromCmd(config.Detection.V6.Cmd)
	} else {
		return false, errors.New("Cannot determine IPv6 because no detection mechanisms are configured")
	}

	if CurrentIPv6 != UpstreamIPv6 {
		DetectionLogger.Info.Printf("Detected change in IPv6 Address: '%s' -> '%s'", CurrentIPv6, UpstreamIPv6)
		CurrentIPv6 = UpstreamIPv6
		return true, nil
	} else {
		DetectionLogger.Debug.Println("No IPv6 change detected")
		return false, nil
	}

}
