/*
	Implementation for desec.io

	Reference:
		- https://desec.readthedocs.io/en/latest/
*/

package services

import (
	"dyngo/config"
	"dyngo/helpers/ip"
	"errors"
	"net/http"
	"strconv"
)

type DesecService struct {
	BaseService
}

func NewDesec() DynDnsService {
	return &DesecService{
		BaseService: NewBaseService("deSEC.io", *config.Services.Desec),
	}
}

func (service *DesecService) Update(Address ip.IPAddress) {
	for _, domain := range service.Domains {
		var err error

		if (domain.V4 && Address.Protocol == ip.IPv4) || (domain.V6 && Address.Protocol == ip.IPv6) {
			err = service.sendUpdateRequest("https://update.dedyn.io", domain.Name, Address.Content)
			service.LogDynDnsUpdate(domain.Name, Address.Content, err)
		}
	}
}

func (service *DesecService) sendUpdateRequest(baseUrl, host, ipAddress string) error {
	var url = baseUrl + "?hostname=" + host + "&myip=" + ipAddress

	service.Logger.Debug.Printf("Sending request: %v\n", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(service.Username, service.Password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected http status: " + strconv.FormatInt(int64(resp.StatusCode), 10))
	} else {
		return nil
	}
}
