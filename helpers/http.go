package helpers

import (
	"dyngo/logger"
	"io"
	"net/http"
)

func ResponseBodyToString(res *http.Response) string {
	body, err := io.ReadAll(res.Body)

	if err != nil {
		logger.Warn.Printf("Could not transform response body to string: %s", err)
		return ""
	}

	return string(body)
}
