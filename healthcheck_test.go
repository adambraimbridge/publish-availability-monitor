package main

import (
	"net/url"
	"testing"
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
