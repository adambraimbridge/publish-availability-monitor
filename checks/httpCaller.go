package checks

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/giantswarm/retry-go"
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
	HTTPMethod, URL, Username, Password, APIKey, TxID, ContentType string
	Entity                                                         io.Reader
}

func NewHttpCaller(timeoutSeconds int) HttpCaller {
	var client http.Client
	if timeoutSeconds > 0 {
		client = http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	return defaultHttpCaller{&client}
}

// Performs http GET calls using the default http client
func (c defaultHttpCaller) DoCall(config Config) (resp *http.Response, err error) {
	if config.HTTPMethod == "" {
		config.HTTPMethod = "GET"
	}
	req, err := http.NewRequest(config.HTTPMethod, config.URL, config.Entity)
	if config.Username != "" && config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}

	if config.APIKey != "" {
		req.Header.Add("X-Api-Key", config.APIKey)
	}

	if config.TxID != "" {
		req.Header.Add("X-Request-Id", config.TxID)
	}

	if config.ContentType != "" {
		req.Header.Add("Content-Type", config.ContentType)
	}

	req.Header.Add("User-Agent", "UPP Publish Availability Monitor")

	op := func() error {
		resp, err = c.client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode >= 500 && resp.StatusCode < 600 {
			//Error status code: create an err in order to trigger a retry
			return fmt.Errorf("Error status code received: %d", resp.StatusCode)
		}
		return nil
	}

	retry.Do(op, retry.RetryChecker(func(err error) bool { return err != nil }), retry.MaxTries(2))
	return resp, err
}
