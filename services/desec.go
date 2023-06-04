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

func (service *DesecService) Update(domain *Domain) error {
	url := "https://update.dedyn.io?hostname=" + domain.Name + "&myip=" + domain.State.Target

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
