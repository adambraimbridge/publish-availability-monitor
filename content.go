package main

import (
	"encoding/json"
	"fmt"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
)

// Content is the interface for different type of contents from different CMSs.
type Content interface {
	isValid() bool
	isMarkedDeleted() bool
	getType() string
	getUUID() string
}

const systemIDKey = "Origin-System-Id"

// Unmarshals the message body into the appropriate content type based on the systemID header.
func unmarshalContent(msg consumer.Message) (Content, error) {
	headers := msg.Headers
	systemID := headers[systemIDKey]
	switch systemID {
	case "http://cmdb.ft.com/systems/methode-web-pub":
		var eomFile EomFile
		err := json.Unmarshal([]byte(msg.Body), &eomFile)
		return eomFile, err
	case "http://cmdb.ft.com/systems/wordpress":
		var wordPressMsg WordPressMessage
		err := json.Unmarshal([]byte(msg.Body), &wordPressMsg)
		return wordPressMsg, err
	default:
		return nil, fmt.Errorf("Unsupported content with system ID: [%s].", systemID)
	}
}
