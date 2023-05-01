package detection

import (
	"dyngo/config"
	"dyngo/detection/strategies"
	"dyngo/helpers/protocols"
	"dyngo/logger"
	"dyngo/services"
	"errors"
)

var DetectionLogger = logger.NewLoggerCollection("detection")
var CurrentIPv4 = ""
var CurrentIPv6 = ""

func RunDetection(trigger string) {
	var IPv4Changed bool = false
	var IPv6Changed bool = false
	var err error

	DetectionLogger.Debug.Printf("Detection triggered (trigger=%s)", trigger)

	if services.AtLeastOneDomainRequires(protocols.IPv4) {
		IPv4Changed, err = DetectIPAddress(protocols.IPv4)

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	}

	if services.AtLeastOneDomainRequires(protocols.IPv6) {
		IPv6Changed, err = DetectIPAddress(protocols.IPv6)

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	}

	if IPv4Changed {
		services.UpdateServices(protocols.IPv4, CurrentIPv4)
	}

	if IPv6Changed {
		services.UpdateServices(protocols.IPv6, CurrentIPv6)
	}
}

func DetectIPAddress(protocol protocols.InternetProtocol) (bool, error) {
	DetectionLogger.Debug.Printf("Running %s detection", protocol.Version)

	var CurrentIPAddress *string
	var ExternalIPAddress string = ""
	var Strategy config.DetectionStrategy

	if protocol == protocols.IPv4 {
		CurrentIPAddress = &CurrentIPv4
		Strategy = config.Detection.Strategies.V4
	} else {
		CurrentIPAddress = &CurrentIPv6
		Strategy = config.Detection.Strategies.V6
	}

	if Strategy.Web != "" {
		ExternalIPAddress = strategies.GetIpAddressFromExternalService(Strategy.Web)
	} else if Strategy.Cmd != "" {
		ExternalIPAddress = strategies.GetIpAddressFromCmd(Strategy.Cmd)
	} else {
		return false, errors.New("Cannot determine " + protocol.Version + " because no detection strategies are configured")
	}

	if *CurrentIPAddress != ExternalIPAddress {
		DetectionLogger.Info.Printf("%s Address changed: '%s' -> '%s' ", protocol.Version, *CurrentIPAddress, ExternalIPAddress)
		*CurrentIPAddress = ExternalIPAddress
		return true, nil
	} else {
		DetectionLogger.Debug.Printf("No %s change detected", protocol.Version)
		return false, nil
	}
}
