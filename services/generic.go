/*
	Implementation for a generic service
*/

package services

import (
	"dyngo/config"
	"dyngo/helpers/dns"
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
	for _, record := range domain.Records {
		var err error

		if service.Protocol == "dyndns2" {
			err = service.useDynDns2Protocol(domain, record)
		} else {
			err = errors.New("Unknown protocol " + service.Protocol)
		}

		service.LogDynDnsUpdate(domain.Name, domain.State[record].Target, err)

		if err != nil {
			return err
		}
	}

	return nil
}

func (service *GenericService) useDynDns2Protocol(domain *Domain, record dns.Record) error {
	state := domain.State[record]

	var url = service.URL + "/nic/update?hostname=" + domain.Name + "&myip=" + state.Target

	service.Logger.Debug.Printf("Sending request: %v\n", url)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(service.Username, service.Password)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	err = service.parseResponse(resp)

	if err != nil {
		return err
	} else {
		domain.HandleSuccessfulUpdate(record)
		return nil
	}
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
