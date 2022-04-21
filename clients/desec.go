package clients

import (
	"dyngo/logger"
	"net/http"
	"strconv"
)

type desec struct {
	service  string
	username string
	password string
}

func NewDesec(username, password string) desec {
	return desec{
		service:  "desec.io",
		username: username,
		password: password,
	}
}

func (d desec) UpdateIPv4(ipAddress, host string) {
	logger.Info.Println("Initiating DNS Update (" + d.service + ")")
	d.sendUpdateRequest("https://update.dedyn.io", host, ipAddress)
}

func (d desec) UpdateIPv6(ipAddress, host string) {
	logger.Info.Println("Initiating DNS Update (" + d.service + ")")
	d.sendUpdateRequest("https://update6.dedyn.io", host, ipAddress)
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
