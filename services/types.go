package services

import (
	"dyngo/config"
	"dyngo/helpers/ip"
	"dyngo/logger"
	"time"
)

type IService interface {
	Update(*Domain) error

	GetName() string
	GetDomains() []*Domain
	GetRetries() int

	// The service might be locked temporarily if subsequent update calls fail.
	// If it is locked, no requests shall be made until it unlocks again.
	IsLocked() bool

	// Resets the internal retry counter. Call this function whenever the services
	// update succeeds.
	ResetRetries()

	// This function increases the internal retry count and returns the number
	// of seconds the service will be locked for. Call this function whenever
	// the services update fails.
	IncreaseRetries() int

	LogDynDnsUpdate(string, string, error)
}

type Domain struct {
	Name     string
	Protocol ip.InternetProtocol
	State    DomainState
}

type DomainState struct {
	Current string
	Target  string
}

type BaseService struct {
	Name     string
	Username string
	Password string
	Domains  []*Domain
	Config   config.ServiceConfiguration
	Logger   *logger.LoggerCollection

	retries    int
	retryAfter time.Time // Timestamp after which the service can be retried
}

func NewBaseService(name string, config config.ServiceConfiguration) BaseService {
	return BaseService{
		Username: config.Username,
		Password: config.Password,
		Domains:  getDomainsFromConfig(config.Domains),
		Name:     name,
		Logger:   logger.NewLoggerCollection("service/" + name),
	}
}

func NewBaseServiceFromGeneric(name string, config config.GenericServiceConfiguration) BaseService {
	return BaseService{
		Username: config.Username,
		Password: config.Password,
		Domains:  getDomainsFromConfig(config.Domains),
		Name:     name,
		Logger:   logger.NewLoggerCollection("service/" + name),
	}
}

func getDomainsFromConfig(domains []config.DomainConfiguration) []*Domain {
	var result []*Domain

	for _, domain := range domains {
		if domain.V4 {
			result = append(result, &Domain{
				Name:     domain.Name,
				Protocol: ip.IPv4,
				State: DomainState{
					Current: "",
					Target:  "",
				},
			})
		}

		if domain.V6 {
			result = append(result, &Domain{
				Name:     domain.Name,
				Protocol: ip.IPv6,
				State: DomainState{
					Current: "",
					Target:  "",
				},
			})
		}
	}

	return result
}

func (service *BaseService) GetDomains() []*Domain {
	return service.Domains
}

func (service *BaseService) GetName() string {
	return service.Name
}

func (service *BaseService) ResetRetries() {
	service.retries = 0
}

func (service *BaseService) GetRetries() int {
	return service.retries
}

func (service *BaseService) IsLocked() bool {
	return time.Now().Compare(service.retryAfter) != 1
}

func (service *BaseService) IncreaseRetries() int {
	service.retries++

	waitInSeconds := 0

	if service.retries < 3 {
		waitInSeconds = 5
	} else if service.retries < 5 {
		waitInSeconds = 60
	} else if service.retries < 10 {
		waitInSeconds = 60 * 10
	} else {
		waitInSeconds = 60 * 30
	}

	service.retryAfter = time.Now().Add(time.Second * time.Duration(waitInSeconds))

	return waitInSeconds
}

func (service *BaseService) LogDynDnsUpdate(domain, ip string, err error) {
	if err == nil {
		service.Logger.Info.Printf("Update '%v' -> '%v' successful", domain, ip)
	} else {
		service.Logger.Error.Printf("Update '%v' -> '%v' failed: %v", domain, ip, err.Error())
	}
}
