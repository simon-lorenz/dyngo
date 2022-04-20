package services

import "dyngo/logger"

var Registered []DynDnsService

func Register(service DynDnsService) {
	Registered = append(Registered, service)
	logger.Info.Printf("Registered service '%v'", service.GetName())
}

func AtLeastOneDomainRequires(protocol string) bool {
	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if (protocol == "v4" && domain.V4) || (protocol == "v6" && domain.V6) {
				return true
			}
		}
	}

	return false
}
