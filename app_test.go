package main

import (
	"os"
	"testing"
)

const syntheticTID = "SYNTHETIC-REQ-MONe4d2885f-1140-400b-9407-921e1c7378cd"
const naturalTID = "tid_xltcnbckvq"

func TestMain(m *testing.M) {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	os.Exit(m.Run())
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