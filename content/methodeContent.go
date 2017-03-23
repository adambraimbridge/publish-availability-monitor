package content

import (
	"strings"

	"encoding/xml"

	"launchpad.net/xmlpath"
)

const sourceXPath = "//ObjectMetadata/EditorialNotes/Sources/Source/SourceCode"
const markDeletedFlagXPath = "//ObjectMetadata/OutputChannels/DIFTcom/DIFTcomMarkDeleted"

// EomFile models Methode content
type EomFile struct {
	UUID             string        `json:"uuid"`
	LinkedObjects    []interface{} `json:"linkedObjects"`
	ContentType      string        `json:"type"`
	Value            string        `json:"value"`
	Attributes       string        `json:"attributes"`
	SystemAttributes string        `json:"systemAttributes"`
	UsageTickets     string        `json:"usageTickets"`
	WorkflowStatus   string        `json:"workflowStatus"`
	Type             string        `json:"-"` //This field is for internal application usage
	Source           Source        `json:"-"` //This field is for internal application usage
	BinaryContent    []byte        `json:"-"` //This field is for internal application usage
}

type Source struct {
	XMLName    xml.Name `xml:"ObjectMetadata"`
	SourceCode string   `xml:"EditorialNotes>Sources>Source>SourceCode"`
}

func (eomfile EomFile) initType() EomFile {
	contentType := eomfile.ContentType
	contentSrc := eomfile.Source.SourceCode

	if contentSrc == "ContentPlaceholder" && contentType == "EOM::CompoundStory" {
		eomfile.Type = "EOM::CompoundStory_ContentPlaceholder"
		infoLogger.Printf("results [%v] ....", eomfile.Type)
		return eomfile
	}
	eomfile.Type = eomfile.ContentType
	return eomfile
}

func (eomfile EomFile) Initialize(binaryContent []byte) Content {
	eomfile.BinaryContent = binaryContent
	return eomfile.initType()
}

func (eomfile EomFile) IsValid(externalValidationEndpoint string, txID string, username string, password string) bool {
	contentUUID := eomfile.UUID
	if !isUUIDValid(contentUUID) {
		warnLogger.Printf("Eomfile invalid: invalid UUID: [%s]. transaction_id=[%s]", contentUUID, txID)
		return false
	}

	return isExternalValidationSuccessful(eomfile.BinaryContent, externalValidationEndpoint, username, password, txID, eomfile.GetUUID(), eomfile.GetType())
}

func (eomfile EomFile) IsMarkedDeleted() bool {
	if eomfile.Type == "Image" || eomfile.Type == "EOM::WebContainer" {
		return false
	}
	markDeletedFlag, ok := getXPathValue(eomfile.Attributes, eomfile, markDeletedFlagXPath)
	if !ok {
		warnLogger.Printf("Eomfile with uuid=[%s]: Cannot match node in XML using xpath [%v]", eomfile.UUID, markDeletedFlagXPath)
		return false
	}
	infoLogger.Printf("Eomfile with uuid=[%s]: MarkAsDeletedFlag: [%v]", eomfile.UUID, markDeletedFlag)
	return markDeletedFlag == "True"
}

func (eomfile EomFile) GetType() string {
	return eomfile.Type
}

func (eomfile EomFile) GetUUID() string {
	return eomfile.UUID
}

func getXPathValue(xml string, eomfile EomFile, lookupPath string) (string, bool) {
	path := xmlpath.MustCompile(lookupPath)
	root, err := xmlpath.Parse(strings.NewReader(xml))
	if err != nil {
		warnLogger.Printf("Cannot parse XML of eomfile with uuid=[%s] using xpath [%v], error: [%v]", eomfile.UUID, lookupPath, err.Error())
		return "", false
	}
	xpathValue, ok := path.String(root)
	return xpathValue, ok

}
