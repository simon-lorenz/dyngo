package services

import "dyngo/logger"

var Registered []DynDnsService

func Register(service DynDnsService) {
	Registered = append(Registered, service)
	logger.Info.Printf("Registered service '%v'", service.GetName())
}
