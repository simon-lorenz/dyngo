package clients

import (
	"dyngo/config"
	"dyngo/logger"
	"net/http"
	"strconv"
)

type desec struct {
	service  string
	username string
	password string
	hosts    []config.HostConfiguration
}

func NewDesec(config config.ServiceConfiguration) desec {
	return desec{
		service:  "desec.io",
		username: config.Username,
		password: config.Password,
		hosts:    config.Hosts,
	}
}

func (d desec) UpdateIPv4(ipAddress string) {
	logger.Info.Println("Initiating DNS Update (" + d.service + ")")

	for _, host := range d.hosts {
		d.sendUpdateRequest("https://update.dedyn.io", host.Host, ipAddress)
	}
}

func (d desec) UpdateIPv6(ipAddress string) {
	logger.Info.Println("Initiating DNS Update (" + d.service + ")")

	for _, host := range d.hosts {
		d.sendUpdateRequest("https://update6.dedyn.io", host.Host, ipAddress)
	}
}

func (d desec) GetHosts() []config.HostConfiguration {
	return d.hosts
}

func (d desec) GetName() string {
	return d.service
}

func (d desec) sendUpdateRequest(baseUrl, host, ipAddress string) {
	var url = baseUrl + "?hostname=" + host + "&myip=" + ipAddress

	logger.Info.Printf("Preparing update request to %v\n", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logger.Error.Println("DNS Update failed: " + err.Error())
		return
	}

	req.SetBasicAuth(d.username, d.password)

	resp, err := client.Do(req)

	if err != nil {
		logger.Error.Println("DNS Update failed: " + err.Error())
	} else if resp.StatusCode != http.StatusOK {
		logger.Error.Println("DNS Update failed: Unexpected http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
	} else {
		logger.Info.Println("DNS Update successful")
	}
}
