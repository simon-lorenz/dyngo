package detection

import (
	"dyngo/config"
	"dyngo/logger"
	"dyngo/services"
	"errors"
)

var DetectionLogger = logger.NewLoggerCollection("detection")
var CurrentIPv4 = ""
var CurrentIPv6 = ""

func RunDetection() {
	var IPv4Changed bool = false
	var IPv6Changed bool = false
	var err error

	if services.AtLeastOneDomainRequires("v4") {
		IPv4Changed, err = DetectIPAddress("v4")

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	}

	if services.AtLeastOneDomainRequires("v6") {
		IPv6Changed, err = DetectIPAddress("v6")

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	}

	for _, service := range services.Registered {
		if IPv4Changed {
			service.UpdateIPv4(CurrentIPv4)
		}

		if IPv6Changed {
			service.UpdateIPv6(CurrentIPv6)
		}
	}
}

func DetectIPAddress(protocol string) (bool, error) {
	DetectionLogger.Debug.Printf("Running IP%s detection", protocol)

	// TODO: Use enum or similar
	if protocol != "v4" && protocol != "v6" {
		return false, errors.New("Can not detect ip address for unknown protocol " + protocol)
	}

	var CurrentIPAddress string = ""
	var ExternalIPAddress string = ""
	var Strategy config.DetectionStrategy

	if protocol == "v4" {
		CurrentIPAddress = CurrentIPv4
		Strategy = config.Detection.Strategies.V4
	} else {
		CurrentIPAddress = CurrentIPv6
		Strategy = config.Detection.Strategies.V6
	}

	if Strategy.Web != "" {
		ExternalIPAddress = getIpAddressFromExternalService(Strategy.Web)
	} else if Strategy.Cmd != "" {
		ExternalIPAddress = getIpAddressFromCmd(Strategy.Cmd)
	} else {
		return false, errors.New("Cannot determine IP" + protocol + " because no detection strategies are configured")
	}

	if CurrentIPAddress != ExternalIPAddress {
		DetectionLogger.Info.Printf("Detected change in IP%s Address: '%s' -> '%s' ", protocol, CurrentIPAddress, ExternalIPAddress)

		if protocol == "v4" {
			CurrentIPv4 = ExternalIPAddress
		} else {
			CurrentIPv6 = ExternalIPAddress
		}

		return true, nil
	} else {
		DetectionLogger.Debug.Printf("No IP%s change detected", protocol)
		return false, nil
	}
}
