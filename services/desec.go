/*
	Implementation for desec.io

	Reference:
		- https://desec.readthedocs.io/en/latest/
*/

package services

import (
	"dyngo/config"
	"dyngo/helpers/dns"
	"errors"
	"net/http"
	"strconv"
)

type DesecService struct {
	BaseService
}

func NewDesec() IService {
	return &DesecService{
		BaseService: NewBaseService("deSEC.io", *config.Services.Desec),
	}
}

func (service *DesecService) Update(domain *Domain) error {
	url := "https://update.dedyn.io?hostname=" + domain.Name

	if domain.Wants(dns.A) {
		url += "&ip=" + domain.State[dns.A].Target
	} else {
		url += "&ip="
	}

	if domain.Wants(dns.AAAA) {
		url += "&ipv6=" + domain.State[dns.AAAA].Target
	} else {
		url += "&ipv6="
	}

	service.Logger.Debug.Printf("Sending request: %v\n", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(service.Username, service.Password)

	resp, err := client.Do(req)

	if err != nil {
		for _, record := range domain.Records {
			service.LogDynDnsUpdate(domain.Name, domain.State[record].Target, err)
		}

		return err
	} else if resp.StatusCode != http.StatusOK {
		err := errors.New("unexpected http status: " + strconv.FormatInt(int64(resp.StatusCode), 10))

		for _, record := range domain.Records {
			service.LogDynDnsUpdate(domain.Name, domain.State[record].Target, err)
		}

		return err
	} else {
		for _, record := range domain.Records {
			service.LogDynDnsUpdate(domain.Name, domain.State[record].Target, err)
			domain.HandleSuccessfulUpdate(record)
		}

		return nil
	}
}
