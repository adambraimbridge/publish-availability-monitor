package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"launchpad.net/xmlpath"
)

var validWordPressTypes []string

func init() {
	validWordPressTypes = []string{
		"post",
		"webchat-markets-live",
		"webchat-live-blogs",
		"webchat-live-qa",
	}
}

// Content is the interface for different type of contents from different CMSs.
type Content interface {
	isValid() bool
	isMarkedDeleted() bool
	getType() string
	getUUID() string
}

// EomFile models Methode content
type EomFile struct {
	UUID             string `json:"uuid"`
	Type             string `json:"type"`
	Value            string `json:"value"`
	Attributes       string `json:"attributes"`
	SystemAttributes string `json:"systemAttributes"`
}

// WordPressMessage models messages from Wordpress
type WordPressMessage struct {
	Status      string `json:"status"`
	APIURL      string `json:"apiUrl"`
	Error       string `json:"error"`
	Post        *Post  `json:"post"`
	PreviousURL string `json:"previousUrl"`
}

// Post models WordPress content
// neglect unused fields (e.g. slug, title, content, etc)
type Post struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

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

func (eomfile EomFile) isValid() bool {
	contentUUID := eomfile.UUID
	if !isUUIDValid(contentUUID) {
		warn.Printf("Eomfile invalid: invalid UUID: [%s]", contentUUID)
		return false
	}

	contentType := eomfile.Type
	switch contentType {
	case webContainer:
		return isListValid(eomfile)
	case compoundStory:
		return isCompoundStoryValid(eomfile)
	case image:
		return isImageValid(eomfile)
	default:
		warn.Printf("Eomfile invalid: unexpected content type: [%s]", contentType)
		return false
	}
}

func (eomfile EomFile) isMarkedDeleted() bool {
	if eomfile.Type == "Image" || eomfile.Type == "EOM::WebContainer" {
		return false
	}
	attributes := eomfile.Attributes
	markDeletedFlagXPath := "//ObjectMetadata/OutputChannels/DIFTcom/DIFTcomMarkDeleted"
	path := xmlpath.MustCompile(markDeletedFlagXPath)
	root, err := xmlpath.Parse(strings.NewReader(attributes))
	if err != nil {
		warn.Printf("Cannot parse attribute XML of eomFile, error: [%v]", err.Error())
		return false
	}
	markDeletedFlag, ok := path.String(root)
	if !ok {
		warn.Printf("Cannot match node in XML using xpath [%v]", markDeletedFlagXPath)
		return false
	}
	info.Printf("MarkAsDeletedFlag: [%v]", markDeletedFlag)
	if markDeletedFlag == "True" {
		return true
	}
	return false

}

func (eomfile EomFile) getType() string {
	return eomfile.Type
}

func (eomfile EomFile) getUUID() string {
	return eomfile.UUID
}

func (wordPressMessage WordPressMessage) isValid() bool {
	contentUUID := wordPressMessage.Post.UUID
	if !isUUIDValid(contentUUID) {
		warn.Printf("WordPress message invalid: invalid UUID: [%s]", contentUUID)
		return false
	}

	apiURL := wordPressMessage.APIURL
	if !isValidBrand(apiURL) {
		warn.Printf("WordPress message invalid: failed to resolve brand for uri [%s].", apiURL)
		return false
	}

	contentType := wordPressMessage.Post.Type
	for _, validType := range validWordPressTypes {
		if contentType == validType {
			return true
		}
	}
	warn.Printf("WordPress message invalid: unexpected content type: [%s]", contentType)
	return false
}

func (wordPressMessage WordPressMessage) isMarkedDeleted() bool {
	if wordPressMessage.Post == nil {
		return true
	}
	return false
}

func (wordPressMessage WordPressMessage) getType() string {
	return wordPressMessage.Post.Type
}

func (wordPressMessage WordPressMessage) getUUID() string {
	return wordPressMessage.Post.UUID
}
