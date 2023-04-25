package detection

import (
	"dyngo/config"
	"dyngo/logger"
	"io"
	"net/http"
	"strconv"
)

func GetIPv4() string {
	return getIpAddressFromExternalService(config.Detection.V4.Web)
}

func GetIPv6() string {
	return getIpAddressFromExternalService(config.Detection.V6.Web)
}

func getIpAddressFromExternalService(url string) string {
	var resp, err = http.Get(url)

	if err != nil {
		logger.Error.Println(err)
		return ""
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		logger.Error.Println("Could not detect ip address via webservice: http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
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
