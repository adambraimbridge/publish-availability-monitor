package main

import (
	"strings"

	"launchpad.net/xmlpath"
)

func isMarkedDeleted(eomFile EomFile) bool {
	if eomFile.Type == "Image" || eomFile.Type == "EOM::WebContainer" {
		return false
	}
	attributes := eomFile.Attributes
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

//TODO extract here all eomfile-util functionality from the validator
