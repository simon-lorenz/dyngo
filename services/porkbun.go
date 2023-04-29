package services

import (
	"bytes"
	"dyngo/config"
	"dyngo/helpers"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type PorkbunService struct {
	BaseService
}

func NewPorkbun() DynDnsService {
	return &PorkbunService{
		BaseService: NewBaseService("porkbun", config.Services.Porkbun),
	}
}

func (service *PorkbunService) UpdateAllDomains(TargetIPv4, TargetIPv6 string) {
	/*
		Reference:
		- https://porkbun.com/api/json/v3/documentation
		- https://github.com/ddclient/ddclient/blob/841ffcbdaa009687e5fb390c4527055e929f959a/ddclient.in#L7199
	*/

	for _, domain := range service.Domains {
		subdomain, host := helpers.ExtractSubdomain((domain.Domain))

		if domain.V4 && TargetIPv4 != "" {
			currentIpAddress, err := service.getExistingRecord("A", host, subdomain)

			if err != nil {
				service.Logger.Info.Println(err)
				continue
			}

			if currentIpAddress == "" {
				service.createRecord(host, subdomain, "A", TargetIPv4)
			} else {
				if currentIpAddress != TargetIPv4 {
					err := service.updateRecord(host, subdomain, "A", TargetIPv4)
					service.LogDynDnsUpdate(domain.Domain, TargetIPv4, err)
				} else {
					service.Logger.Info.Printf("Current ip address for %s does not differ from target ip address, skipping", domain.Domain)
				}
			}
		}

		if domain.V6 && TargetIPv6 != "" {
			currentIpAddress, err := service.getExistingRecord("AAAA", host, subdomain)

			if err != nil {
				service.Logger.Info.Println(err)
				continue
			}

			if currentIpAddress == "" {
				service.createRecord(host, subdomain, "AAAA", TargetIPv6)
			} else {
				if currentIpAddress != TargetIPv6 {
					err := service.updateRecord(host, subdomain, "AAAA", TargetIPv6)
					service.LogDynDnsUpdate(domain.Domain, TargetIPv6, err)
				} else {
					service.Logger.Info.Printf("Current ip address for %s does not differ from target ip address, skipping", domain.Domain)
				}
			}
		}
	}
}

func (service *PorkbunService) getExistingRecord(record, domain, subdomain string) (string, error) {
	var ENDPOINT string = "https://porkbun.com/api/json/v3/dns/retrieveByNameType/" + domain + "/" + record + "/" + subdomain

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

func (service *PorkbunService) createRecord(domain, subdomain, record, ip string) (string, error) {
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
		"type":         record,
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

func (service *PorkbunService) updateRecord(domain, subdomain, record, ip string) error {
	service.Logger.Info.Printf("Updating %s record for %s", record, helpers.JoinDomainParts(subdomain, domain))

	var ENDPOINT string = "https://porkbun.com/api/json/v3/dns/editByNameType/" + domain + "/" + record + "/" + subdomain

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
