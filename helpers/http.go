package helpers

import (
	"io"
	"net/http"
)

func ResponseBodyToString(res *http.Response) string {
	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return string(body)
}
