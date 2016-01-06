package content

import (
	"testing"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
)

func TestUnmarshalContent_ValidMessageMethodeSystemHeader_NoError(t *testing.T) {
	if _, err := UnmarshalContent(validMethodeMessage); err != nil {
		t.Errorf("Message with valid system ID [%s] cannot be unmarshalled!", validMethodeMessage.Headers["Origin-System-Id"])
	}
}

func TestUnmarshalContent_ValidMessageWordpressSystemHeader_NoError(t *testing.T) {
	if _, err := UnmarshalContent(validWordpressMessage); err != nil {
		t.Errorf("Message with valid system ID [%s] cannot be unmarshalled!", validWordpressMessage.Headers["Origin-System-Id"])
	}
}

func TestUnmarshalContent_InvalidMessageMissingHeader_Error(t *testing.T) {
	if _, err := UnmarshalContent(invalidMessageWrongHeader); err == nil {
		t.Error("Expected failure, but message with missing system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_InvalidMessageWrongSystemId_Error(t *testing.T) {
	if _, err := UnmarshalContent(invalidMessageWrongSystemID); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_InvalidMethodeContentWrongJSONFormat_Error(t *testing.T) {
	if _, err := UnmarshalContent(invalidMethodeMessageWrongJSONFormat); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_InvalidWordPressContentWrongJSONFormat_Error(t *testing.T) {
	if _, err := UnmarshalContent(invalidWordPressMessageWrongJSONFormat); err == nil {
		t.Error("Expected failure, but message with wrong system ID successfully unmarshalled!")
	}
}

func TestUnmarshalContent_ValidWordPressMessageWithTypeField_TypeIsCorrectlyUnmarshalled(t *testing.T) {
	content, err := UnmarshalContent(validWordPressMessageWithTypeField)
	if err != nil {
		t.Errorf("Expected success, but error occured [%v]", err)
		return
	}
	if content.GetType() != "post" {
		t.Errorf("Expected [post] content type, but found [%s].", content.GetType())
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
