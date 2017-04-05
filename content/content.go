package content

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/publish-availability-monitor/checks"
)

// Content is the interface for different type of contents from different CMSs.
type Content interface {
	Validate(externalValidationEndpoint string, txID string, username string, password string) ValidationResponse
	GetType() string
	GetUUID() string
}

type ValidationResponse struct {
	IsValid         bool
	IsMarkedDeleted bool
}

const systemIDKey = "Origin-System-Id"

var httpCaller checks.HttpCaller

func init() {
	httpCaller = checks.NewHttpCaller(10)
}

// UnmarshalContent unmarshals the message body into the appropriate content type based on the systemID header.
func UnmarshalContent(msg consumer.Message) (Content, error) {
	headers := msg.Headers

	systemID := headers[systemIDKey]
	switch systemID {
	case "http://cmdb.ft.com/systems/methode-web-pub":
		var eomFile EomFile

		err := json.Unmarshal([]byte(msg.Body), &eomFile)
		if err != nil {
			return nil, err
		}
		xml.Unmarshal([]byte(eomFile.Attributes), &eomFile.Source)
		eomFile = eomFile.initType()
		return eomFile, err
	case "http://cmdb.ft.com/systems/wordpress":
		var wordPressMsg WordPressMessage
		err := json.Unmarshal([]byte(msg.Body), &wordPressMsg)
		return wordPressMsg, err
	case "http://cmdb.ft.com/systems/brightcove":
		var video Video
		err := json.Unmarshal([]byte(msg.Body), &video)
		return video, err
	default:
		return nil, fmt.Errorf("Unsupported content with system ID: [%s].", systemID)
	}
}
