package checks

import (
	"io"
	"net/http"
	"time"
)

// httpCaller abstracts http calls
type HttpCaller interface {
	DoCall(url string, username string, password string, txId string) (*http.Response, error)
	DoCallWithEntity(httpMethod string, url string, username string, password string, txId string, contentType string, entity io.Reader) (*http.Response, error)
}

// Default implementation of httpCaller
type defaultHttpCaller struct {
	client *http.Client
}

func NewHttpCaller(timeoutSeconds int) HttpCaller {
	var client http.Client
	if timeoutSeconds > 0 {
		client = http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	} else {
		client = http.Client{}
	}

	return defaultHttpCaller{&client}
}

// Performs http GET calls using the default http client
func (c defaultHttpCaller) DoCall(url string, username string, password string, txId string) (resp *http.Response, err error) {
	return c.DoCallWithEntity("GET", url, username, password, txId, "", nil)
}

func (c defaultHttpCaller) DoCallWithEntity(httpMethod string, url string, username string, password string, txId string, contentType string, entity io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(httpMethod, url, entity)
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if txId != "" {
		req.Header.Add("X-Request-Id", txId)
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	req.Header.Add("User-Agent", "UPP Publish Availability Monitor")

	return c.client.Do(req)
}
