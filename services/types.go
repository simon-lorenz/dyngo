package services

import (
	"dyngo/config"
	"dyngo/helpers"
	"dyngo/helpers/dns"
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
	Name    string
	State   map[dns.Record]DomainState
	Records []dns.Record
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

func (domain *Domain) UpdateTarget(address ip.IPAddress) {
	var record dns.Record

	if address.Protocol == ip.IPv4 {
		record = dns.A
	} else if address.Protocol == ip.IPv6 {
		record = dns.AAAA
	} else {
		panic("Cannot update target for record '" + string(record) + "'")
	}

	if state, ok := domain.State[record]; ok {
		state.Target = address.Content
		domain.State[record] = state
	}
}

func (domain *Domain) NeedsUpdate() bool {
	for _, record := range domain.Records {
		if domain.State[record].Current != domain.State[record].Target {
			return true
		}
	}

	return false
}

func (domain *Domain) Wants(record dns.Record) bool {
	return helpers.Contains(domain.Records, record)
}

func (domain *Domain) HandleSuccessfulUpdate(record dns.Record) {
	state := domain.State[record]
	state.Current = state.Target
	domain.State[record] = state
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
	return helpers.Map(domains, func(config config.DomainConfiguration) *Domain {
		domain := Domain{
			Name:  config.Name,
			State: make(map[dns.Record]DomainState),
		}

		for _, record := range config.Records {
			domain.Records = append(domain.Records, record)

			domain.State[record] = DomainState{
				Current: "",
				Target:  "",
			}
		}

		return &domain
	})
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
