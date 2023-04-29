package detection

import (
	"dyngo/config"
	"dyngo/logger"
	"errors"
)

var CurrentIPv4 = ""
var CurrentIPv6 = ""

func RefreshIPv4() (bool, error) {
	logger.Debug.Println("Running IPv4 detection")

	if config.Detection.V4.Web != "" {
		UpstreamIPv4 := getIpAddressFromExternalService(config.Detection.V4.Web)

		if CurrentIPv4 != UpstreamIPv4 {
			logger.Info.Printf("Detected change in IPv4 Address: '%v' -> '%v' \n", CurrentIPv4, UpstreamIPv4)
			CurrentIPv4 = UpstreamIPv4
			return true, nil
		} else {
			logger.Debug.Println("No IPv4 change detected")
			return false, nil
		}

	} else {
		return false, errors.New("Cannot determine IPv4 because no detection mechanisms are configured.")
	}
}

func RefreshIPv6() (bool, error) {
	logger.Debug.Println("Running IPv6 detection")

	if config.Detection.V6.Web != "" {
		UpstreamIPv6 := getIpAddressFromExternalService(config.Detection.V4.Web)

		if CurrentIPv4 != UpstreamIPv6 {
			logger.Info.Printf("Detected change in IPv6 Address: '%v' -> '%v' \n", CurrentIPv4, UpstreamIPv6)
			CurrentIPv6 = UpstreamIPv6
			return true, nil
		} else {
			logger.Debug.Println("No IPv6 change detected")
			return false, nil
		}
	} else {
		return false, errors.New("Cannot determine IPv6 because no detection mechanisms are configured.")
	}
}
