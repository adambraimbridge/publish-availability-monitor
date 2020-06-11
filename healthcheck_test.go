package main

import (
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildFtHealthcheckURL(t *testing.T) {
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
		if actual, _ := buildFtHealthcheckURL(*uri, tc.health); actual != tc.expectedHealthURL {
			t.Errorf("For [%s]:\n\tExpected: [%s]\n\tActual: [%s]", tc.validationURL, tc.expectedHealthURL, actual)
		}
	}
}

func TestBuildAwsHealthcheckURL(t *testing.T) {
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
		if actual, _ := buildAwsHealthcheckURL(tc.validationURL); actual != tc.expectedHealthURL {
			t.Errorf("For [%s]:\n\tExpected: [%s]\n\tActual: [%s]", tc.validationURL, tc.expectedHealthURL, actual)
		}

	}

}

func TestPublishNoFailuresForSameUUIDs(t *testing.T) {
	config := MetricConfig{}
	interval := Interval{5, 5}
	newURL := url.URL{}
	t0 := time.Now()
	publishMetric1 := PublishMetric{"1234567", false, t0, "", interval, config, newURL, "tid_1234", false}
	publishMetric2 := PublishMetric{"1234567", false, t0, "", interval, config, newURL, "tid_6789", false}
	publishMetric3 := PublishMetric{"1234567", false, t0, "", interval, config, newURL, "tid_6789", false}
	testMetrics := []PublishMetric{publishMetric1, publishMetric2, publishMetric3}
	testPublishHistory := publishHistory{sync.RWMutex{}, testMetrics}
	testHealthcheck := Healthcheck{
		config:          &AppConfig{},
		metricContainer: &testPublishHistory,
	}
	_, err := testHealthcheck.checkForPublishFailures()

	assert.NoError(t, err, "No Error expected if multiple fails for the same uuid")
}

func TestPublishFailureForDistinctUUIDs(t *testing.T) {
	config := MetricConfig{}
	interval := Interval{5, 5}
	newURL := url.URL{}
	t0 := time.Now()
	publishMetric1 := PublishMetric{"12345", false, t0, "", interval, config, newURL, "tid_1234", false}
	publishMetric2 := PublishMetric{"12678", false, t0, "", interval, config, newURL, "tid_6789", false}
	publishMetric3 := PublishMetric{"12679", true, t0, "", interval, config, newURL, "tid_6789", false}
	testMetrics := []PublishMetric{publishMetric1, publishMetric2, publishMetric3}
	testPublishHistory := publishHistory{sync.RWMutex{}, testMetrics}
	testHealthcheck := Healthcheck{
		config:          &AppConfig{},
		metricContainer: &testPublishHistory,
	}
	_, err := testHealthcheck.checkForPublishFailures()

	assert.Error(t, err, "Expected Error for at least two distinct uuid publish fails")
}
