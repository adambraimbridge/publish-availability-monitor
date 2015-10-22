package main

import (
	"testing"

	"github.com/Financial-Times/go-message-queue-consumer"
)

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

func TestIsMessageValid_MessageValid(t *testing.T) {
	if !isMessageValid(validMessage) {
		t.Error("Valid Message marked as invalid!")
	}
}

func TestIsMessageValid_MissingHeader(t *testing.T) {
	if isMessageValid(invalidMessageWrongHeader) {
		t.Error("Invalid Message marked as valid!")
	}
}

func TestIsMessageValid_InvalidSystemId(t *testing.T) {
	if isMessageValid(invalidMessageWrongHeader) {
		t.Error("Invalid Message marked as valid!")
	}
}

var validMessage = consumer.Message{
	map[string]string{
		"Origin-System-Id": "methode-web-pub",
	},
	"body",
}
var invalidMessageWrongHeader = consumer.Message{
	map[string]string{
		"Foobar-System-Id": "methode-web-pub",
	},
	"body",
}
var invalidMessageWrongSystemId = consumer.Message{
	map[string]string{
		"Origin-System-Id": "methode-web-foobar",
	},
	"body",
}
var validUUID = "e28b12f7-9796-3331-b030-05082f0b8157"
var invalidUUID = "foobar"
