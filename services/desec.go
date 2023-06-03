/*
	Implementation for desec.io

	Reference:
		- https://desec.readthedocs.io/en/latest/
*/

package services

import (
	"dyngo/config"
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

func (service *DesecService) Update() error {
	for _, domain := range service.GetDomains() {
		if domain.State.Current == domain.State.Target {
			continue
		}

		err := service.sendUpdateRequest("https://update.dedyn.io", domain.Name, domain.State.Target)
		service.LogDynDnsUpdate(domain.Name, domain.State.Target, err)

		if err != nil {
			return err
		}

		domain.State.Current = domain.State.Target
	}

	return nil
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
