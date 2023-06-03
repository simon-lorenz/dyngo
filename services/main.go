package services

import (
	"dyngo/helpers/ip"
	"dyngo/logger"
	"strconv"
	"time"
)

var ServiceLogger = logger.NewLoggerCollection("service")

var Registered []IService

func Register(service IService) {
	Registered = append(Registered, service)
	ServiceLogger.Info.Printf("Registered service '%v'", service.GetName())
}

func SetTarget(IPAddress ip.IPAddress) {
	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if domain.Protocol == IPAddress.Protocol {
				domain.State.Target = IPAddress.Content
			}
		}
	}
}

func GetServicesThatNeedUpdate() []IService {
	var NeedUpdate []IService

	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if domain.State.Current != domain.State.Target {
				NeedUpdate = append(NeedUpdate, service)
				break
			}
		}
	}

	return NeedUpdate
}

func AtLeastOneDomainRequires(protocol ip.InternetProtocol) bool {
	for _, service := range Registered {
		for _, domain := range service.GetDomains() {
			if domain.Protocol == protocol {
				return true
			}
		}
	}

	return false
}

func init() {
	go func() {
		for range time.Tick(time.Second * 1) {
			for _, service := range GetServicesThatNeedUpdate() {
				if service.IsLocked() {
					// This service shouldn't be retried this tick.
					continue
				}

				err := service.Update()

				if err == nil {
					service.ResetRetries()
				} else {
					lockedFor := service.IncreaseRetries()
					ServiceLogger.Warn.Printf("Error while updating service %q. Will retry in %ss.", service.GetName(), strconv.Itoa(lockedFor))
				}
			}
		}
	}()
}
