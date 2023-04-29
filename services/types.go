package services

import (
	"dyngo/config"
	"dyngo/logger"
)

type DynDnsService interface {
	UpdateIPv4(string)
	UpdateIPv6(string)

	GetDomains() []config.DomainConfiguration
	GetName() string
}

type BaseService struct {
	Name     string
	Username string
	Password string
	Domains  []config.DomainConfiguration
	Config   config.ServiceConfiguration
	Logger   *logger.LoggerCollection
}

func NewBaseService(name string, config config.ServiceConfiguration) BaseService {
	return BaseService{
		Username: config.Username,
		Password: config.Password,
		Domains:  config.Domains,
		Name:     name,
		Logger:   logger.NewLoggerCollection(name),
	}
}

func (service *BaseService) GetDomains() []config.DomainConfiguration {
	return service.Domains
}

func (service *BaseService) GetName() string {
	return service.Name
}

func (service *BaseService) LogDynDnsUpdate(domain, ip string, err error) {
	if err == nil {
		service.Logger.Info.Printf("Update '%v' -> '%v' successful", domain, ip)
	} else {
		service.Logger.Error.Printf("Update '%v' -> '%v' failed: %v", domain, ip, err.Error())
	}
}
