package strategies

import (
	"io"
	"net/http"
	"strconv"
)

type WebStrategy struct {
	BaseDetectionStrategy
	URL string
}

func NewWebDetectionStrategy(URL string) IDetectionStrategy {
	return &WebStrategy{
		BaseDetectionStrategy: NewBaseDetectionStrategy("web"),
		URL:                   URL,
	}
}

func (strategy *WebStrategy) Execute() string {
	strategy.Logger.Debug.Printf("URL: %s", strategy.URL)

	var resp, err = http.Get(strategy.URL)

	if err != nil {
		strategy.Logger.Error.Println(err)
		return ""
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		strategy.Logger.Error.Println("Could not detect ip address: http status " + strconv.FormatInt(int64(resp.StatusCode), 10))
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		strategy.Logger.Error.Printf("Could not detect ip address: %s", err)
		return ""
	}

	ip := string(body)

	strategy.Logger.Debug.Printf("Detection successful: %s", ip)

	return ip
}
