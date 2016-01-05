package main

import (
	"testing"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
)

func TestUnmarshalContent_ValidMessageMethodeSystemHeader_NoError(t *testing.T) {
	if _, err := unmarshalContent(validMethodeMessage); err != nil {
		t.Errorf("Message with valid system ID [%s] cannot be unmarshalled!", validMethodeMessage.Headers["Origin-System-Id"])
	}
}

func TestUnmarshalContent_ValidMessageWordpressSystemHeader_NoError(t *testing.T) {
	if _, err := unmarshalContent(validWordpressMessage); err != nil {
		t.Errorf("Message with valid system ID [%s] cannot be unmarshalled!", validWordpressMessage.Headers["Origin-System-Id"])
	}
}

func TestUnmarshalContent_InvalidMessageMissingHeader_Error(t *testing.T) {
	if _, err := unmarshalContent(invalidMessageWrongHeader); err == nil {
		t.Error("Expected failure, but message with missing system ID successfully unmarshalled!")
	}
}

func TestIsMessageValid_InvalidSystemId_Error(t *testing.T) {
	if _, err := unmarshalContent(invalidMessageWrongSystemID); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

var validMethodeMessage = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
	},
	Body: "{}",
}
var validWordpressMessage = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/wordpress",
	},
	Body: "{}",
}
var invalidMessageWrongHeader = consumer.Message{
	Headers: map[string]string{
		"Foobar-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
	},
	Body: "{}",
}
var invalidMessageWrongSystemID = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "methode-web-foobar",
	},
	Body: "{}",
}
