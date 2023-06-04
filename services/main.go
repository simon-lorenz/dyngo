package services

import (
	"dyngo/helpers"
	"dyngo/helpers/dns"
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
			domain.UpdateTarget(IPAddress)
		}
	}
}

func GetServicesThatNeedUpdate() []IService {
	return helpers.Filter(Registered, func(service IService) bool {
		return helpers.Find(service.GetDomains(), func(domain *Domain) bool {
			return domain.NeedsUpdate()
		}) != nil
	})
}

func AtLeastOneDomainRequires(protocol ip.InternetProtocol) bool {
	return helpers.Find(Registered, func(service IService) bool {
		return helpers.Find(service.GetDomains(), func(domain *Domain) bool {
			return domain.Wants(dns.GetRecordForInternetProtocol(protocol))
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
					if !domain.NeedsUpdate() {
						continue
					}

					err := service.Update(domain)

					if err == nil {
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
