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

func TestUnmarshalContent_InvalidMessageWrongSystemId_Error(t *testing.T) {
	if _, err := unmarshalContent(invalidMessageWrongSystemID); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_InvalidMethodeContentWrongJSONFormat_Error(t *testing.T) {
	if _, err := unmarshalContent(invalidMethodeMessageWrongJSONFormat); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_InvalidWordPressContentWrongJSONFormat_Error(t *testing.T) {
	if _, err := unmarshalContent(invalidWordPressMessageWrongJSONFormat); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_ValidWordPressMessageWithTypeField_TypeIsCorrectlyUnmarshalled(t *testing.T) {
	content, err := unmarshalContent(validWordPressMessageWithTypeField)
	if err != nil {
		t.Errorf("Expected success, but error occured [%v]", err)
		return
	}
	if content.getType() != "post" {
		t.Errorf("Expected [post] content type, but found [%s].", content.getType())
	}
}

var invalidMethodeMessageWrongJSONFormat = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
	},
	Body: `{"uuid": "79e7f5ed-63c7-46b2-9767-736f8ae3a3f6", "type": "Image", "value" : " }`,
}

var validMethodeMessage = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/methode-web-pub",
	},
	Body: "{}",
}

var invalidWordPressMessageWrongJSONFormat = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/wordpress",
	},
	Body: `{"status": "ok", "post": {"id : "002251", "type": "post"}}`,
}

var validWordPressMessageWithTypeField = consumer.Message{
	Headers: map[string]string{
		"Origin-System-Id": "http://cmdb.ft.com/systems/wordpress",
	},
	Body: `{"status": "ok", "post": {"id" : "002251", "type": "post"}}`,
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
