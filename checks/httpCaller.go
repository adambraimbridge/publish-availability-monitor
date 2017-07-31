package checks

import (
	"io"
	"net/http"
	"time"
)

// httpCaller abstracts http calls
type HttpCaller interface {
	DoCall(config Config) (*http.Response, error)
}

// Default implementation of httpCaller
type defaultHttpCaller struct {
	client *http.Client
}

type Config struct {
	HttpMethod, Url, Username, Password, ApiKey, TxId, ContentType string
	Entity                                                         io.Reader
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
func (c defaultHttpCaller) DoCall(config Config) (resp *http.Response, err error) {
	if config.HttpMethod == "" {
		config.HttpMethod = "GET"
	}
	req, err := http.NewRequest(config.HttpMethod, config.Url, config.Entity)
	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	if config.ApiKey != "" {
		req.Header.Add("X-Api-Key", config.ApiKey)
	}

	if config.TxId != "" {
		req.Header.Add("X-Request-Id", config.TxId)
	}

	if config.ContentType != "" {
		req.Header.Add("Content-Type", config.ContentType)
	}

	req.Header.Add("User-Agent", "UPP Publish Availability Monitor")

	return c.client.Do(req)
}
