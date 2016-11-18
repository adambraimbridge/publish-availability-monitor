package checks

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthenticated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "UPP Publish Availability Monitor", r.Header.Get("User-Agent"), "user agent header")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world"))
	}))
	defer server.Close()

	httpCaller := NewHttpCaller(10)
	resp, err := httpCaller.DoCall(server.URL, "", "")
	assert.Nil(t, err, "unexpected error")

	by := make([]byte, 32)
	i, _ := resp.Body.Read(by)
	body := string(by[:i])

	assert.Equal(t, http.StatusOK, resp.StatusCode, "response status")
	assert.Equal(t, "Hello world", body, "response body")
}

func TestAuthenticated(t *testing.T) {
	username := "scott"
	password := "tiger"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "UPP Publish Availability Monitor", r.Header.Get("User-Agent"), "user agent header")

		expectedAuthz := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		assert.Equal(t, expectedAuthz, r.Header.Get("Authorization"), "authorization header")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello world"))
	}))
	defer server.Close()

	httpCaller := NewHttpCaller(10)
	resp, err := httpCaller.DoCall(server.URL, username, password)
	assert.Nil(t, err, "unexpected error")

	by := make([]byte, 32)
	i, _ := resp.Body.Read(by)
	body := string(by[:i])

	assert.Equal(t, http.StatusOK, resp.StatusCode, "response status")
	assert.Equal(t, "Hello world", body, "response body")
}
