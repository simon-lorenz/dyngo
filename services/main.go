package services

import (
	"dyngo/helpers"
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
	return helpers.Filter(Registered, func(service IService) bool {
		return helpers.Find(service.GetDomains(), func(domain *Domain) bool {
			return domain.State.Current != domain.State.Target
		}) != nil
	})
}

func AtLeastOneDomainRequires(protocol ip.InternetProtocol) bool {
	return helpers.Find(Registered, func(service IService) bool {
		return helpers.Find(service.GetDomains(), func(domain *Domain) bool {
			return domain.Protocol == protocol
		}) != nil
	}) != nil
}

func init() {
	go func() {
		for range time.Tick(time.Second * 1) {
			for _, service := range GetServicesThatNeedUpdate() {
				if service.IsLocked() {
					// This service shouldn't be retried this tick.
					continue
				}

				for _, domain := range service.GetDomains() {
					if domain.State.Current == domain.State.Target {
						continue
					}

					err := service.Update(domain)
					service.LogDynDnsUpdate(domain.Name, domain.State.Target, err)

					if err == nil {
						domain.State.Current = domain.State.Target
						service.ResetRetries()
					} else {
						lockedFor := service.IncreaseRetries()
						ServiceLogger.Warn.Printf("Error while updating service %q. Will retry in %ss.", service.GetName(), strconv.Itoa(lockedFor))
						break
					}
				}

			}
		}
	}()
}
