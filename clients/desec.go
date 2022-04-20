package clients

import (
	"dyngo/config"
	"dyngo/logger"
	"net/http"
	"strconv"
	"strings"
)

func DesecIPv4(ipAddress string) {
	logger.Info.Println("Initiating DNS Update (Desec.io)")
	sendUpdateRequest("https://update.dedyn.io", ipAddress)
}

func DesecIPv6(ipAddress string) {
	logger.Info.Println("Initiating DNS Update (Desec.io)")
	sendUpdateRequest("https://update6.dedyn.io", ipAddress)
}

func sendUpdateRequest(baseUrl, ipAddress string) {
	var url = baseUrl + "?hostname=" + strings.Join(config.Domains, ",") + "&myip=" + ipAddress

	logger.Info.Printf("Preparing update request to %v\n", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		logger.Error.Println("DNS Update failed: " + err.Error())
		return
	}

	req.SetBasicAuth(config.Username, config.Token)

	resp, err := client.Do(req)

	if err != nil {
		logger.Error.Println("DNS Update failed: " + err.Error())
	} else if resp.StatusCode != http.StatusOK {
		logger.Error.Println("DNS Update failed: Unexpected http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
	} else {
		logger.Info.Println("DNS Update successful")
	}
}
