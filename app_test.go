package main

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	os.Exit(m.Run())
}

func TestIsMessagePastPublishSLA_pastSLA(t *testing.T) {
	publishDate := time.Now().Add(-(threshold + 1) * time.Second)
	if !isMessagePastPublishSLA(publishDate, threshold) {
		t.Error("Did not detect message past SLA")
	}
}

func TestIsMessagePastPublishSLA_notPastSLA(t *testing.T) {
	publishDate := time.Now()
	if isMessagePastPublishSLA(publishDate, threshold) {
		t.Error("Valid message marked as passed SLA")
	}
}

func TestIsIgnorableMessage_naturalMessage(t *testing.T) {
	if isIgnorableMessage(naturalTID) {
		t.Error("Normal message marked as ignorable")
	}
}

func TestIsIgnorableMessage_syntheticMessage(t *testing.T) {
	if !isIgnorableMessage(syntheticTID) {
		t.Error("Synthetic message marked as normal")
	}
}

func TestGetCredentials(t *testing.T) {
	environments["env1"] = Environment{"env1", "http://env1.example.org", "http://s3.example.org", "user1", "pass1"}
	environments["env2"] = Environment{"env2", "http://env2.example.org", "http://s3.example.org", "user2", "pass2"}

	username, password := getCredentials("http://env2.example.org/__some-service")
	if username != "user2" || password != "pass2" {
		t.Error("incorrect credentials returned")
	}
}

func TestGetCredentials_Unauthenticated(t *testing.T) {
	environments["env1"] = Environment{"env1", "http://env1.example.org", "http://s3.example.org", "user1", "pass1"}
	environments["env2"] = Environment{"env2", "http://env2.example.org", "http://s3.example.org", "user2", "pass2"}

	username, password := getCredentials("http://env3.example.org/__some-service")
	if username != "" || password != "" {
		t.Error("incorrect credentials returned")
	}
}

const threshold = 120
const syntheticTID = "SYNTHETIC-REQ-MONe4d2885f-1140-400b-9407-921e1c7378cd"
const naturalTID = "tid_xltcnbckvq"
