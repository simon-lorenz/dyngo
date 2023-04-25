package detection

import (
	"dyngo/config"
	"errors"
)

func GetIPv4() (string, error) {
	if config.Detection.V4.Web != "" {
		return getIpAddressFromExternalService(config.Detection.V4.Web), nil
	} else {
		return "", errors.New("Cannot determine IPv4 because no detection mechanisms are configured.")
	}
}

func GetIPv6() (string, error) {
	if config.Detection.V6.Web != "" {
		return getIpAddressFromExternalService(config.Detection.V6.Web), nil
	} else {
		return "", errors.New("Cannot determine IPv6 because no detection mechanisms are configured.")
	}
}
