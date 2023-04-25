package services

import (
	"dyngo/config"
	"dyngo/logger"
	"errors"
	"net/http"
	"strconv"
)

type desec struct {
	username   string
	password   string
	targetIPv4 string
	targetIPv6 string
	domains    []DynDnsDomain
}

func NewDesec(config config.ServiceConfiguration) DynDnsService {
	var result desec = desec{
		username: config.Username,
		password: config.Password,
	}

	for _, domain := range config.Domains {
		result.domains = append(result.domains, DynDnsDomain{
			domain:      domain.Domain,
			V4:          domain.V4,
			V6:          domain.V6,
			currentIpV4: "",
			currentIPv6: "",
		})
	}

	return &result
}

func (service *desec) GetDomains() []DynDnsDomain {
	return service.domains
}

func (service *desec) GetName() string {
	return "deSEC.io"
}

func (service *desec) SetTargetIPv4(ip string) {
	service.targetIPv4 = ip
}

func (service *desec) SetTargetIPv6(ip string) {
	service.targetIPv6 = ip
}

func (service *desec) UpdateAllDomains() {
	for i := range service.domains {
		domain := &service.domains[i]

		if domain.V4 && domain.currentIpV4 != service.targetIPv4 {
			err := service.sendUpdateRequest("https://update.dedyn.io", domain.domain, service.targetIPv4)

			if err == nil {
				logger.LogDynDnsUpdate(service.GetName(), domain.domain, service.targetIPv4, nil)
				domain.currentIpV4 = service.targetIPv4
			} else {
				logger.LogDynDnsUpdate(service.GetName(), domain.domain, service.targetIPv6, err)
			}
		}

		if domain.V6 && domain.currentIPv6 != service.targetIPv6 {
			err := service.sendUpdateRequest("https://update6.dedyn.io", domain.domain, service.targetIPv6)

			if err == nil {
				logger.LogDynDnsUpdate(service.GetName(), domain.domain, service.targetIPv4, nil)
				domain.currentIPv6 = service.targetIPv6
			} else {
				logger.LogDynDnsUpdate(service.GetName(), domain.domain, service.targetIPv6, err)
			}
		}
	}
}

func (service *desec) sendUpdateRequest(baseUrl, host, ipAddress string) error {
	var url = baseUrl + "?hostname=" + host + "&myip=" + ipAddress

	logger.Debug.Printf("[%v] Sending request: %v\n", service.GetName(), url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(service.username, service.password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected http status: " + strconv.FormatInt(int64(resp.StatusCode), 10))
	} else {
		return nil
	}
}
