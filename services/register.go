package services

import (
	"dyngo/helpers/protocols"
	"dyngo/logger"
)

var Registered []DynDnsService

func Register(service DynDnsService) {
	Registered = append(Registered, service)
	logger.Info.Printf("Registered service '%v'", service.GetName())
}

func UpdateServices(protocol protocols.InternetProtocol, IPAddress string) {
	for _, service := range Registered {
		if protocol == protocols.IPv4 {
			service.UpdateIPv4(IPAddress)
		}

		if protocol == protocols.IPv6 {
			service.UpdateIPv6(IPAddress)
		}
	}
}

func AtLeastOneDomainRequires(protocol protocols.InternetProtocol) bool {
	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if (protocol == protocols.IPv4 && domain.V4) || (protocol == protocols.IPv6 && domain.V6) {
				return true
			}
		}
	}

	return false
}
