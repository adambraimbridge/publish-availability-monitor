package main

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/Financial-Times/go-message-queue-consumer"
	"github.com/satori/go.uuid"
	"launchpad.net/xmlpath"
)

const image = "Image"
const webContainer = "EOM::WebContainer"
const compoundStory = "EOM::CompoundStory"
const systemIdKey = "Origin-System-Id"

const titleXPath = "/doc/lead/lead-headline/headline/ln"
const channelXPath = "/props/productInfo/name"
const webTypeXPath = "//ObjectMetadata/FTcom/DIFTcomWebType"
const filePathXPath = "//ObjectMetadata/EditorialNotes/ObjectLocation"

const expectedSystemId = "methode-web-pub"
const expectedWebChannel = "FTcom"
const expectedWebTypePrefix = "digitalList"
const expectedFilePathSuffix = ".xml"

func isMessageValid(message consumer.Message) bool {
	headers := message.Headers
	systemId := headers[systemIdKey]
	if systemId != expectedSystemId {
		//TODO message
		return false
	}
	return true
}

func isEomfileValid(eomfile EomFile) bool {
	contentUuid := eomfile.UUID
	parsedUuid, _ := uuid.FromString(contentUuid)
	//TODO handle error

	if contentUuid != parsedUuid.String() {
		//TODO message
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
		//TODO message
		return false
	}
}

func isListValid(eomfile EomFile) bool {
	attributes := eomfile.Attributes
	path := xmlpath.MustCompile(webTypeXPath)
	root, err := xmlpath.Parse(strings.NewReader(attributes))
	if err != nil {
		log.Fatal(err)
	}
	webType, _ := path.String(root)
	//TODO handle not ok
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
		log.Fatal(err)
	}
	filePath, _ := path.String(root)
	//TODO handle not ok
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
		log.Fatal(err)
	}
	channel, _ := path.String(root)
	//TODO handle not ok
	if channel == expectedWebChannel {
		return true
	}
	return false
}

func hasTitle(eomfile EomFile) bool {
	if eomfile.Value == nil || len(eomfile.Value) == 0 {
		//TODO message
		return false
	}
	//decode value from base64
	decoded, err := base64.StdEncoding.DecodeString(string(eomfile.Value[:]))
	if err != nil {
		log.Printf("ERROR - failure in decoding base64 value: %s", err.Error())
		return false
	}
	articleXml := string(decoded[:])

	path := xmlpath.MustCompile(titleXPath)
	root, err := xmlpath.Parse(strings.NewReader(articleXml))
	if err != nil {
		log.Fatal(err)
	}
	title, _ := path.String(root)
	//TODO handle not ok
	if len(title) > 0 {
		return true
	}
	return false
}

func isImageValid(eomfile EomFile) bool {
	if eomfile.Value == nil || len(eomfile.Value) == 0 {
		//TODO message
		return false
	}
	return true
}
