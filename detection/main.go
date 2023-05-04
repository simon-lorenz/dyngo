package detection

import (
	"dyngo/config"
	"dyngo/detection/strategies"
	"dyngo/helpers/ip"
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

	if services.AtLeastOneDomainRequires(ip.IPv4) {
		IPv4Changed, err = DetectIPAddress(ip.IPv4)

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	} else {
		DetectionLogger.Debug.Println("Skipping IPv4 detection because no service needs it.")
	}

	if services.AtLeastOneDomainRequires(ip.IPv6) {
		IPv6Changed, err = DetectIPAddress(ip.IPv6)

		if err != nil {
			DetectionLogger.Error.Println(err.Error())
		}
	} else {
		DetectionLogger.Debug.Println("Skipping IPv6 detection because no service needs it.")
	}

	if IPv4Changed {
		services.UpdateServices(ip.IPAddress{
			Content:  CurrentIPv4,
			Protocol: ip.IPv4,
		})
	}

	if IPv6Changed {
		services.UpdateServices(ip.IPAddress{
			Content:  CurrentIPv6,
			Protocol: ip.IPv6,
		})
	}
}

func DetectIPAddress(protocol ip.InternetProtocol) (bool, error) {
	DetectionLogger.Debug.Printf("Running %s detection", protocol.Version)

	var CurrentIPAddress *string
	var AvailableStrategies config.DetectionStrategies
	var Strategy strategies.DetectionStrategy

	if protocol == ip.IPv4 {
		CurrentIPAddress = &CurrentIPv4
		AvailableStrategies = config.Detection.Strategies.V4
	} else {
		CurrentIPAddress = &CurrentIPv6
		AvailableStrategies = config.Detection.Strategies.V6
	}

	if AvailableStrategies.Web != "" {
		Strategy = strategies.NewWebDetectionStrategy(AvailableStrategies.Web)
	} else if AvailableStrategies.Cmd != "" {
		Strategy = strategies.NewCmdDetectionStrategy(AvailableStrategies.Cmd)
	} else {
		return false, errors.New("Cannot determine " + protocol.Version + " because no detection strategies are configured")
	}

	ExternalIPAddress := Strategy.Execute()

	if *CurrentIPAddress != ExternalIPAddress {
		DetectionLogger.Info.Printf("%s Address changed: '%s' -> '%s' ", protocol.Version, *CurrentIPAddress, ExternalIPAddress)
		*CurrentIPAddress = ExternalIPAddress
		return true, nil
	} else {
		DetectionLogger.Debug.Printf("No %s change detected", protocol.Version)
		return false, nil
	}
}
