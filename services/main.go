package services

import (
	"dyngo/helpers/ip"
	"dyngo/logger"
)

var Registered []DynDnsService

func Register(service DynDnsService) {
	Registered = append(Registered, service)
	logger.Info.Printf("Registered service '%v'", service.GetName())
}

func UpdateServices(IPAddress ip.IPAddress) {
	for _, service := range Registered {
		service.Update(IPAddress)
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
