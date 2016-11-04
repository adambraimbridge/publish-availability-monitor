package checks

import (
	"net/http"
	"time"
)

// httpCaller abstracts http calls
type HttpCaller interface {
	DoCall(url string, username string, password string) (*http.Response, error)
}

// Default implementation of httpCaller
type defaultHttpCaller struct {
	client *http.Client
}

func NewHttpCaller(timeoutSeconds int) HttpCaller {
	var client http.Client
	if timeoutSeconds > 0 {
		client = http.Client{Timeout: time.Duration(10 * time.Second)}
	} else {
		client = http.Client{}
	}

	return defaultHttpCaller{&client}
}

// Performs http GET calls using the default http client
func (c defaultHttpCaller) DoCall(url string, username string, password string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	req.Header.Add("User-Agent", "UPP Publish Availability Monitor")

	return c.client.Do(req)
}
