package strategies

import (
	"dyngo/logger"
	"io"
	"net/http"
	"strconv"
)

var webDetectionLogger = logger.NewLoggerCollection("detection/strategies/web")

func GetIpAddressFromExternalService(url string) string {
	webDetectionLogger.Debug.Printf("URL: %s", url)

	var resp, err = http.Get(url)

	if err != nil {
		webDetectionLogger.Error.Println(err)
		return ""
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		webDetectionLogger.Error.Println("Could not detect ip address: http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		webDetectionLogger.Error.Printf("Could not detect ip address: %s", err)
		return ""
	}

	ip := string(body)

	webDetectionLogger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
