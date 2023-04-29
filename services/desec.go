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

func NewDesec() DynDnsService {
	return &DesecService{
		BaseService: NewBaseService("deSEC.io", *config.Services.Desec)}
}

func (service *DesecService) UpdateIPv4(Target string) {
	if Target == "" {
		return
	}

	for i := range service.Domains {
		domain := &service.Domains[i]

		if domain.V4 {
			err := service.sendUpdateRequest("https://update.dedyn.io", domain.Domain, Target)

			if err == nil {
				service.LogDynDnsUpdate(domain.Domain, Target, nil)
			} else {
				service.LogDynDnsUpdate(domain.Domain, Target, err)
			}
		}
	}
}

func (service *DesecService) UpdateIPv6(Target string) {
	if Target == "" {
		return
	}

	for i := range service.Domains {
		domain := &service.Domains[i]

		if domain.V6 {
			err := service.sendUpdateRequest("https://update.dedyn.io", domain.Domain, Target)

			if err == nil {
				service.LogDynDnsUpdate(domain.Domain, Target, nil)
			} else {
				service.LogDynDnsUpdate(domain.Domain, Target, err)
			}
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
