package main

import (
	"strings"

	"launchpad.net/xmlpath"
)

// EomFile models Methode content
type EomFile struct {
	UUID             string `json:"uuid"`
	Type             string `json:"type"`
	Value            string `json:"value"`
	Attributes       string `json:"attributes"`
	SystemAttributes string `json:"systemAttributes"`
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
