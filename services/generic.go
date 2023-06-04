/*
	Implementation for a generic service
*/

package services

import (
	"dyngo/config"
	"errors"
	"io"
	"net/http"
	"strings"
)

type GenericService struct {
	Protocol string
	URL      string
	BaseService
}

func NewGenericService(config config.GenericServiceConfiguration) IService {
	return &GenericService{
		Protocol:    config.Protocol,
		URL:         config.URL,
		BaseService: NewBaseServiceFromGeneric("generic/"+strings.ToLower(config.Name), config),
	}
}

func (service *GenericService) Update(domain *Domain) error {
	if service.Protocol == "dyndns2" {
		return service.useDynDns2Protocol(domain.Name, domain.State.Target)
	} else {
		return errors.New("Unknown protocol " + service.Protocol)
	}
}

func (service *GenericService) useDynDns2Protocol(host, ipAddress string) error {
	var url = service.URL + "/nic/update?hostname=" + host + "&myip=" + ipAddress

	service.Logger.Debug.Printf("Sending request: %v\n", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(service.Username, service.Password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	return service.parseResponse(resp)
}

func (service *GenericService) parseResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return errors.New("Could not parse response")
	}

	code := string(body)

	if code == "good" {
		return nil
	}

	if code == "nochg" {
		// This is more a warning than an error. It can happen but should be prevented
		// whenever possible, otherwise the client might get blocked.
		return errors.New("Received return code nochg")
	}

	fatalReturnCodes := []string{"abuse", "badagent", "badauth", "badsys", "dnserr", "!donator", "nohost", "notfqdn", "numhost", "!yours", "911"}

	for _, fatalReturnCode := range fatalReturnCodes {
		if code == fatalReturnCode {
			return errors.New("Received fatal return code: " + code)
		}
	}

	return nil
}
