package services

import (
	"dyngo/helpers/ip"
	"dyngo/logger"
	"strconv"
	"time"
)

var ServiceLogger = logger.NewLoggerCollection("service")
var Registered []DynDnsService

func Register(service DynDnsService) {
	Registered = append(Registered, service)
	ServiceLogger.Info.Printf("Registered service '%v'", service.GetName())
}

func UpdateServices(IPAddress ip.IPAddress) {
	const MAX_RETRIES = 2

	for _, service := range Registered {
		ServiceLogger.Debug.Printf("Updating %s for service %q", IPAddress.Protocol.Version, service.GetName())

		for i := 0; i <= MAX_RETRIES; i++ {
			err := service.Update(IPAddress)

			if err != nil {
				if i != MAX_RETRIES {
					ServiceLogger.Warn.Printf("Error while updating service %q. Will retry in 5s (%s/%s)", service.GetName(), strconv.Itoa(i+1), strconv.Itoa(MAX_RETRIES))
					time.Sleep(time.Second * 5)
				} else {
					ServiceLogger.Warn.Printf("Failed to update service %q. Will wait for next trigger.", service.GetName())
				}
			} else {
				break
			}
		}
	}
}

func AtLeastOneDomainRequires(protocol ip.InternetProtocol) bool {
	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if (protocol == ip.IPv4 && domain.V4) || (protocol == ip.IPv6 && domain.V6) {
				return true
			}
		}
	}

	return false
}
