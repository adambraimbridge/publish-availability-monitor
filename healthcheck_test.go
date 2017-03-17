package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestBuildFtHealthcheckUrl(t *testing.T) {
	var testCases = []struct {
		validationURL     string
		health            string
		expectedHealthURL string
	}{
		{
			validationURL:     "http://methode-article-transformer/content-transformer/",
			health:            "/__methode-article-transformer/__health",
			expectedHealthURL: "http://methode-article-transformer/__methode-article-transformer/__health",
		},
		{
			validationURL:     "http://methode-article-transformer/content-transformer?monitor=true",
			health:            "/__methode-article-transformer/__health",
			expectedHealthURL: "http://methode-article-transformer/__methode-article-transformer/__health",
		},
	}
	for _, tc := range testCases {
		uri, _ := url.Parse(tc.validationURL)
		if actual, _ := buildFtHealthcheckUrl(*uri, tc.health); actual != tc.expectedHealthURL {
			t.Errorf("For [%s]:\n\tExpected: [%s]\n\tActual: [%s]", tc.validationURL, tc.expectedHealthURL, actual)
		}
	}
}

func TestBuildAwsHealthcheckUrl(t *testing.T) {
	var testCases = []struct {
		validationURL     string
		expectedHealthURL string
	}{
		{
			validationURL:     "http://some-bucket.amazonaws.com/",
			expectedHealthURL: "http://some-bucket.amazonaws.com/healthCheckDummyFile",
		},
	}
	for _, tc := range testCases {
		if actual, _ := buildAwsHealthcheckUrl(tc.validationURL); actual != tc.expectedHealthURL {
			t.Errorf("For [%s]:\n\tExpected: [%s]\n\tActual: [%s]", tc.validationURL, tc.expectedHealthURL, actual)
		}

	}

}

func TestPublishNoFailuresForSametUUIDs(t *testing.T) {
	assert := assert.New(t)
	config := MetricConfig{}
	interval := Interval{5, 5}
	newUrl := url.URL{}
	t0 := time.Now()
	var publishMetric1 = PublishMetric{"1234567", false, t0, "", interval, config, newUrl, "tid_1234", false}
	var publishMetric2 = PublishMetric{"1234567", false, t0, "", interval, config, newUrl, "tid_6789", false}
	var publishMetric3 = PublishMetric{"1234567", false, t0, "", interval, config, newUrl, "tid_6789", false}
	var testMetrics = []PublishMetric{publishMetric1, publishMetric2, publishMetric3}

	var testPublishHistory = publishHistory{sync.RWMutex{}, testMetrics}
	var testHealthcheck = Healthcheck{http.Client{}, AppConfig{}, &testPublishHistory}
	_, err := testHealthcheck.checkForPublishFailures()

	assert.NoError(err, "No Error expected if multiple fails for the same uuid")

}

func TestPublishFailureForDistinctUUIDs(t *testing.T) {
	assert := assert.New(t)
	config := MetricConfig{}
	interval := Interval{5, 5}
	newUrl := url.URL{}
	t0 := time.Now()
	var publishMetric1 = PublishMetric{"12345", false, t0, "", interval, config, newUrl, "tid_1234", false}
	var publishMetric2 = PublishMetric{"12678", false, t0, "", interval, config, newUrl, "tid_6789", false}
	var publishMetric3 = PublishMetric{"12679", true, t0, "", interval, config, newUrl, "tid_6789", false}
	var testMetrics = []PublishMetric{publishMetric1, publishMetric2, publishMetric3}

	testPublishHistory := publishHistory{sync.RWMutex{}, testMetrics}

	testHealthcheck := Healthcheck{http.Client{}, AppConfig{}, &testPublishHistory}
	_, err := testHealthcheck.checkForPublishFailures()

	assert.Error(err, "Expected Error for at least two distinct uuid publish fails")
}
