package checks

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stubServer(t *testing.T, expectedHeaders map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, v := range expectedHeaders {
			actualValue := r.Header.Get(k)
			assert.Equal(t, v, actualValue, k+" header")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world"))
	}))

	return server
}

func assertExpectedResponse(t *testing.T, resp *http.Response) {
	by := make([]byte, 32)
	i, _ := resp.Body.Read(by)
	body := string(by[:i])

	assert.Equal(t, http.StatusOK, resp.StatusCode, "response status")
	assert.Equal(t, "Hello world", body, "response body")
}

func TestUnauthenticated(t *testing.T) {
	server := stubServer(t, map[string]string{
		"User-Agent": "UPP Publish Availability Monitor",
	})
	defer server.Close()

	httpCaller := NewHttpCaller(10)
	resp, err := httpCaller.DoCall(server.URL, "", "", "")
	assert.Nil(t, err, "unexpected error")

	assertExpectedResponse(t, resp)
}

func TestAuthenticated(t *testing.T) {
	username := "scott"
	password := "tiger"

	server := stubServer(t, map[string]string{
		"User-Agent":    "UPP Publish Availability Monitor",
		"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password)),
	})
	defer server.Close()

	httpCaller := NewHttpCaller(10)
	resp, err := httpCaller.DoCall(server.URL, username, password, "")
	assert.Nil(t, err, "unexpected error")

	assertExpectedResponse(t, resp)
}

func TestTransactionId(t *testing.T) {
	txId := "tid_myTxId"

	server := stubServer(t, map[string]string{
		"User-Agent":   "UPP Publish Availability Monitor",
		"X-Request-Id": txId,
	})
	defer server.Close()

	httpCaller := NewHttpCaller(10)
	resp, err := httpCaller.DoCall(server.URL, "", "", txId)
	assert.Nil(t, err, "unexpected error")

	assertExpectedResponse(t, resp)
}
