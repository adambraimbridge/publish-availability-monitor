package main

import "testing"

func TestBuildFtHealthcheckUrl(t *testing.T) {
	var testCases = []struct {
		validationURL     string
		expectedHealthURL string
	}{
		{
			validationURL:     "http://methode-article-transformer/content-transformer/",
			expectedHealthURL: "http://methode-article-transformer/__health",
		},
		{
			validationURL:     "http://methode-article-transformer/content-transformer",
			expectedHealthURL: "http://methode-article-transformer/__health",
		},
		{
			validationURL:     "http://methode-article-transformer/__methode-article-transformer/content-transformer",
			expectedHealthURL: "http://methode-article-transformer/__methode-article-transformer/__health",
		},
		{
			validationURL:     "http://methode-article-transformer",
			expectedHealthURL: "http://methode-article-transformer/__health",
		},
		{
			validationURL:     "http://methode-article-transformer:8080",
			expectedHealthURL: "http://methode-article-transformer:8080/__health",
		},
		{
			validationURL:     "http://localhost:8080/__methode-article-transformer/content-transformer/",
			expectedHealthURL: "http://localhost:8080/__methode-article-transformer/__health",
		},
		{
			validationURL:     "http://coco.example.com/__notifications-rw/content/notifications",
			expectedHealthURL: "http://coco.example.com/__notifications-rw/__health",
		},
	}
	for _, tc := range testCases {
		if actual, _ := buildFtHealthcheckUrl(tc.validationURL); actual != tc.expectedHealthURL {
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
