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

func TestIsSyntheticMessage_naturalMessage(t *testing.T) {
	if isSyntheticMessage(naturalTID) {
		t.Error("Normal message marked as synthetic")
	}
}

func TestIsSyntheticMessage_syntheticMessage(t *testing.T) {
	if !isSyntheticMessage(syntheticTID) {
		t.Error("Synthetic message marked as normal")
	}
}

func TestIsUUIDValid_UUIDValid(t *testing.T) {
	if !isUUIDValid(validUUID) {
		t.Error("Valid UUID marked as invalid!")
	}
}

func TestIsUUIDValid_UUIDInvalid(t *testing.T) {
	if isUUIDValid(invalidUUID) {
		t.Error("Invalid UUID marked as valid!")
	}
}

const validUUID = "e28b12f7-9796-3331-b030-05082f0b8157"
const invalidUUID = "foobar"
const threshold = 120
const syntheticTID = "SYNTHETIC-REQ-MONe4d2885f-1140-400b-9407-921e1c7378cd"
const naturalTID = "tid_xltcnbckvq"
