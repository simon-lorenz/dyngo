package detection

import (
	"dyngo/config"
	"dyngo/logger"
	"io"
	"net/http"
)

func GetIPv4() string {
	return getIpAddressFromExternalService(config.IPv4AddressDetection.Web)
}

func GetIPv6() string {
	return getIpAddressFromExternalService(config.IPv6AddressDetection.Web)
}

func getIpAddressFromExternalService(url string) string {
	var resp, err = http.Get(url)

	if err != nil || resp.StatusCode < 200 || resp.StatusCode > 300 {
		logger.Error.Println(err)
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Error.Println(err)
		return ""
	}

	return string(body)
}
