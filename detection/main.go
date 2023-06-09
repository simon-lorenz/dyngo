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
		services.SetTarget(ip.IPAddress{
			Content:  CurrentIPv4,
			Protocol: ip.IPv4,
		})
	}

	if IPv6Changed {
		services.SetTarget(ip.IPAddress{
			Content:  CurrentIPv6,
			Protocol: ip.IPv6,
		})
	}
}

func DetectIPAddress(protocol ip.InternetProtocol) (bool, error) {
	DetectionLogger.Debug.Printf("Running %s detection", protocol.Version)

	var CurrentIPAddress *string

	if protocol == ip.IPv4 {
		CurrentIPAddress = &CurrentIPv4
	} else {
		CurrentIPAddress = &CurrentIPv6
	}

	AvailableStrategies := GetAvailableDetectionStrategies(protocol)

	if len(AvailableStrategies) == 0 {
		return false, errors.New("Cannot determine " + protocol.Version + " because no detection strategies are configured")
	}

	ExternalIPAddress := ""

	for _, Strategy := range AvailableStrategies {
		ExternalIPAddress = Strategy.Execute()

		if ExternalIPAddress != "" {
			break
		}
	}

	if ExternalIPAddress == "" {
		return false, errors.New("Could not determine " + protocol.Version + " because no detection strategy succeeded")
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

func GetAvailableDetectionStrategies(protocol ip.InternetProtocol) []strategies.IDetectionStrategy {
	var StrategiesForProtocol config.DetectionStrategies

	if protocol == ip.IPv4 {
		StrategiesForProtocol = config.Detection.Strategies.V4
	} else {
		StrategiesForProtocol = config.Detection.Strategies.V6
	}

	AvailableStrategies := make([]strategies.IDetectionStrategy, 0)

	if StrategiesForProtocol.Web != "" {
		AvailableStrategies = append(AvailableStrategies, strategies.NewWebDetectionStrategy(StrategiesForProtocol.Web))
	}

	if StrategiesForProtocol.Cmd != "" {
		AvailableStrategies = append(AvailableStrategies, strategies.NewCmdDetectionStrategy(StrategiesForProtocol.Cmd))
	}

	return AvailableStrategies
}
