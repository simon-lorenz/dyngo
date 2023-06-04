/*
	Implementation for porkbun.com

	Reference:
		- https://porkbun.com/api/json/v3/documentation
		- https://github.com/ddclient/ddclient/blob/841ffcbdaa009687e5fb390c4527055e929f959a/ddclient.in#L7199
*/

package services

import (
	"bytes"
	"dyngo/config"
	"dyngo/helpers"
	"dyngo/helpers/dns"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type PorkbunService struct {
	BaseService
}

func NewPorkbun() IService {
	return &PorkbunService{
		BaseService: NewBaseService("porkbun", *config.Services.Porkbun),
	}
}

func (service *PorkbunService) Update(domain *Domain) error {
	for _, record := range domain.Records {
		targetIpAddress := domain.State[record].Target
		subdomain, host := helpers.ExtractSubdomain((domain.Name))
		recordIPAddress, err := service.getExistingRecord(record, host, subdomain)

		if err != nil {
			service.LogDynDnsUpdate(domain.Name, targetIpAddress, err)
			continue
		}

		if recordIPAddress == "" {
			_, err = service.createRecord(host, subdomain, record, targetIpAddress)
		} else {
			if recordIPAddress != targetIpAddress {
				err = service.updateRecord(host, subdomain, record, targetIpAddress)
			} else {
				service.Logger.Info.Printf("Current ip address for %s does not differ from target ip address, skipping", domain.Name)
			}
		}

		service.LogDynDnsUpdate(domain.Name, targetIpAddress, err)

		if err != nil {
			return err
		} else {
			domain.HandleSuccessfulUpdate(record)
		}
	}

	return nil
}

func (service *PorkbunService) getExistingRecord(record dns.Record, domain, subdomain string) (string, error) {
	var ENDPOINT string = "https://porkbun.com/api/json/v3/dns/retrieveByNameType/" + domain + "/" + string(record) + "/" + subdomain

	type Response struct {
		Status     string `json:"status"`
		Cloudflare string `json:"cloudflare"`
		Records    []struct {
			Id      string `json:"id"`
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
			TTL     string `json:"ttl"`
			Prio    string `json:"prio"`
			Notes   string `json:"notes"`
		} `json:"records"`
	}

	body, _ := json.Marshal(map[string]interface{}{
		"apikey":       service.Username,
		"secretapikey": service.Password,
	})

	req, _ := http.NewRequest(http.MethodPost, ENDPOINT, bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		service.Logger.Debug.Printf(helpers.ResponseBodyToString(res))
		return "", errors.New("Unexpected http status: " + strconv.FormatInt(int64(res.StatusCode), 10))
	}

	defer res.Body.Close()

	response := &Response{}
	json.NewDecoder(res.Body).Decode(response)

	for _, r := range response.Records {
		if r.Name == subdomain+"."+domain {
			return r.Content, nil
		}
	}

	service.Logger.Info.Printf("No record found for %s", strings.Join([]string{subdomain, domain}, "."))

	return "", nil
}

func (service *PorkbunService) createRecord(domain, subdomain string, record dns.Record, ip string) (string, error) {
	service.Logger.Info.Printf("Creating new record for %s", helpers.JoinDomainParts(subdomain, domain))

	var ENDPOINT string = "https://porkbun.com/api/json/v3/dns/create/" + domain

	type Response struct {
		Status string `json:"status"`
		Id     string `json:"id"`
	}

	body, _ := json.Marshal(map[string]interface{}{
		"apikey":       service.Username,
		"secretapikey": service.Password,
		"name":         subdomain,
		"type":         string(record),
		"content":      ip,
	})

	req, _ := http.NewRequest(http.MethodPost, ENDPOINT, bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		service.Logger.Debug.Println(helpers.ResponseBodyToString(res))
		return "", errors.New("Unexpected http status: " + strconv.FormatInt(int64(res.StatusCode), 10))
	}

	defer res.Body.Close()

	response := &Response{}
	json.NewDecoder(res.Body).Decode(response)

	return response.Id, nil
}

func (service *PorkbunService) updateRecord(domain, subdomain string, record dns.Record, ip string) error {
	service.Logger.Info.Printf("Updating %s record for %s", string(record), helpers.JoinDomainParts(subdomain, domain))

	var ENDPOINT string = "https://porkbun.com/api/json/v3/dns/editByNameType/" + domain + "/" + string(record) + "/" + subdomain

	body, _ := json.Marshal(map[string]interface{}{
		"apikey":       service.Username,
		"secretapikey": service.Password,
		"content":      ip,
	})

	req, _ := http.NewRequest(http.MethodPost, ENDPOINT, bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		service.Logger.Debug.Printf(helpers.ResponseBodyToString(res))
		return errors.New("Unexpected http status: " + strconv.FormatInt(int64(res.StatusCode), 10))
	}

	return nil
}
