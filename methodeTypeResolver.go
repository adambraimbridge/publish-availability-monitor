package main

import (
	"fmt"
	"encoding/xml"
	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/Financial-Times/publish-availability-monitor/content"
	log "github.com/Sirupsen/logrus"
)

var blogCategories = []string{"blog", "webchat-live-blogs", "webchat-live-qa", "webchat-markets-live", "fastft"}

type typeResolver interface {
	ResolveTypeAndUuid(eomFile content.EomFile, txID string) (string, string, error)
}

type methodeTypeResolver struct{
	iResolver checks.IResolver
}

func NewMethodeTypeResolver(iResolver checks.IResolver) *methodeTypeResolver {
	return &methodeTypeResolver{
		iResolver: iResolver,
	}
}

func (m *methodeTypeResolver) ResolveTypeAndUuid(eomFile content.EomFile, txID string) (string, string, error) {
	contentType := eomFile.ContentType
	contentSrc := eomFile.Source.SourceCode
	if contentSrc == "ContentPlaceholder" && contentType == "EOM::CompoundStory" {
		resolvedUuid, err := m.resolveUUID(eomFile, txID)
		if err != nil {
			return "", "", err
		}

		theType := "EOM::CompoundStory_External_CPH"
		cphUuid := eomFile.UUID
		if resolvedUuid != "" {
			theType = "EOM::CompoundStory_Internal_CPH"
			cphUuid = resolvedUuid
		}
		log.Infof("For placeholder resolved type=%v uuid=%v", theType, cphUuid)
		return theType, cphUuid, nil
	}
	return eomFile.ContentType, eomFile.UUID, nil
}

func (m *methodeTypeResolver) resolveUUID(eomFile content.EomFile, txID string) (string, error) {
	attributes, err := m.buildAttributes(eomFile.Attributes)
	if err != nil {
		return "", err
	}

	uuid := ""
	if isBlogCategory(attributes) {
		resolvedUuid, err := m.iResolver.ResolveIdentifier(attributes.ServiceId, attributes.RefField, txID)
		if err != nil {
			return "", fmt.Errorf("Couldn't resolve blog uuid %v", err)
		}
		uuid = resolvedUuid
	}
	return uuid, nil
}

func (m *methodeTypeResolver) buildAttributes(attributesXML string) (content.Attributes, error) {
	var attrs content.Attributes
	if err := xml.Unmarshal([]byte(attributesXML), &attrs); err != nil {
		return content.Attributes{}, err
	}
	return attrs, nil
}

func isBlogCategory(attributes content.Attributes) bool {
	for _, c := range blogCategories {
		if c == attributes.Category {
			return true
		}
	}
	return false
}
