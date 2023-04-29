package detection

import (
	"dyngo/logger"
	"io"
	"net/http"
	"strconv"
)

var WebDetectionLogger = logger.NewLoggerCollection("detection/web")

func getIpAddressFromExternalService(url string) string {
	WebDetectionLogger.Debug.Printf("URL: %s", url)

	var resp, err = http.Get(url)

	if err != nil {
		WebDetectionLogger.Error.Println(err)
		return ""
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		WebDetectionLogger.Error.Println("Could not detect ip address: http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		WebDetectionLogger.Error.Printf("Could not detect ip address: %s", err)
		return ""
	}

	ip := string(body)

	WebDetectionLogger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
