package main

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/satori/go.uuid"
	"launchpad.net/xmlpath"
)

const image = "Image"
const systemIDKey = "Origin-System-Id"
const webContainer = "EOM::WebContainer"
const compoundStory = "EOM::CompoundStory"

const titleXPath = "/doc/lead/lead-headline/headline/ln"
const channelXPath = "/props/productInfo/name"
const webTypeXPath = "//ObjectMetadata/FTcom/DIFTcomWebType"
const filePathXPath = "//ObjectMetadata/EditorialNotes/ObjectLocation"

const expectedSystemID = "http://cmdb.ft.com/systems/methode-web-pub"
const expectedWebChannel = "FTcom"
const expectedWebTypePrefix = "digitalList"
const expectedFilePathSuffix = ".xml"

const syntheticPrefix = "SYNTHETIC"

func isMessagePastPublishSLA(date time.Time, threshold int) bool {
	passedSLA := date.Add(time.Duration(threshold) * time.Second)
	return time.Now().After(passedSLA)
}

func isSyntheticMessage(tid string) bool {
	return strings.HasPrefix(tid, syntheticPrefix)
}

func isMessageValid(message consumer.Message) bool {
	headers := message.Headers
	systemID := headers[systemIDKey]
	if systemID != expectedSystemID {
		warn.Printf("Message invalid: unexpected system ID: [%v]", systemID)
		return false
	}
	return true
}

func isEomfileValid(eomfile EomFile) bool {
	contentUUID := eomfile.UUID
	if !isUUIDValid(contentUUID) {
		warn.Printf("Eomfile invalid: invalid UUID: [%v]", contentUUID)
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
		warn.Printf("Eomfile invalid: unexpected content type: [%v]", contentType)
		return false
	}
}

func isUUIDValid(contentUUID string) bool {
	parsedUUID, err := uuid.FromString(contentUUID)
	if err != nil {
		warn.Printf("Cannot parse UUID [%v], error: [%v]", contentUUID, err.Error())
		return false
	}
	return contentUUID == parsedUUID.String()
}

func isListValid(eomfile EomFile) bool {
	attributes := eomfile.Attributes
	path := xmlpath.MustCompile(webTypeXPath)
	root, err := xmlpath.Parse(strings.NewReader(attributes))
	if err != nil {
		warn.Printf("Cannot parse attribute XML of eomfile, error: [%v]", err.Error())
		return false
	}
	webType, ok := path.String(root)
	if !ok {
		warn.Printf("Cannot match node in XML using xpath [%v]", webTypeXPath)
		return false
	}
	if strings.HasPrefix(webType, expectedWebTypePrefix) {
		return true
	}
	return false
}

func isCompoundStoryValid(eomfile EomFile) bool {
	return isSupportedFileType(eomfile) && isWebChannel(eomfile) && hasTitle(eomfile)
}

func isSupportedFileType(eomfile EomFile) bool {
	attributes := eomfile.Attributes
	path := xmlpath.MustCompile(filePathXPath)
	root, err := xmlpath.Parse(strings.NewReader(attributes))
	if err != nil {
		warn.Printf("Cannot parse attribute XML of eomfile, error: [%v]", err.Error())
		return false
	}
	filePath, ok := path.String(root)
	if !ok {
		warn.Printf("Cannot match node in XML using xpath [%v]", filePathXPath)
		return false
	}
	if strings.HasSuffix(filePath, expectedFilePathSuffix) {
		return true
	}
	return false
}

func isWebChannel(eomfile EomFile) bool {
	systemAttributes := eomfile.SystemAttributes
	path := xmlpath.MustCompile(channelXPath)
	root, err := xmlpath.Parse(strings.NewReader(systemAttributes))
	if err != nil {
		warn.Printf("Cannot parse system attribute XML of eomfile, error: [%v]", err.Error())
		return false
	}
	channel, ok := path.String(root)
	if !ok {
		warn.Printf("Cannot match node in XML using xpath [%v]", channelXPath)
		return false
	}
	if channel == expectedWebChannel {
		return true
	}
	return false
}

func hasTitle(eomfile EomFile) bool {
	if len(eomfile.Value) == 0 {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(eomfile.Value)
	if err != nil {
		warn.Printf("Cannot decode Base64-encoded eomfile value: [%v]", err.Error())
		return false
	}
	articleXML := string(decoded[:])

	root, err := xmlpath.Parse(strings.NewReader(articleXML))
	if err != nil {
		warn.Printf("Cannot parse value XML of eomfile, error: [%v]", err.Error())
		return false
	}

	path := xmlpath.MustCompile(titleXPath)
	title, ok := path.String(root)
	if !ok {
		warn.Printf("Cannot match node in XML using xpath [%v]", titleXPath)
		return false
	}

	title = strings.TrimSpace(title)
	if len(title) > 0 {
		return true
	}
	return false
}

func isImageValid(eomfile EomFile) bool {
	if len(eomfile.Value) == 0 {
		return false
	}
	return true
}
