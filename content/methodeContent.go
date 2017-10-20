package content

import (
	"encoding/xml"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"fmt"
	"github.com/Financial-Times/publish-availability-monitor/checks"
	"io/ioutil"
	"os"
	"encoding/json"
)

var blogCategories = []string{"blog", "webchat-live-blogs", "webchat-live-qa", "webchat-markets-live", "fastft"}

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
	BinaryContent    []byte        `json:"-"` //This field is for internal application usag
}

type Source struct {
	XMLName    xml.Name `xml:"ObjectMetadata"`
	SourceCode string   `xml:"EditorialNotes>Sources>Source>SourceCode"`
}

// Attributes is the data structure that models methode content placeholders attributes
type Attributes struct {
	XMLName             xml.Name `xml:"ObjectMetadata"`
	SourceCode          string   `xml:"EditorialNotes>Sources>Source>SourceCode"`
	LastPublicationDate string   `xml:"OutputChannels>DIFTcom>DIFTcomLastPublication"`
	RefField            string   `xml:"WiresIndexing>ref_field"`
	ServiceId           string   `xml:"WiresIndexing>serviceid"`
	Category            string   `xml:"WiresIndexing>category"`
	IsDeleted           bool     `xml:"OutputChannels>DIFTcom>DIFTcomMarkDeleted"`
}

var iResolver checks.IResolver

func init() {
	httpCaller = checks.NewHttpCaller(10)
}

func InitializeUUIDResolver(uuidResolverUrl string, readEnvUsername string, readEnvPassword string) {
	docStoreClient := checks.NewHttpDocStoreClient(uuidResolverUrl, httpCaller, readEnvUsername, readEnvPassword)
	iResolver = checks.NewHttpIResolver(docStoreClient, readBrandMappings())
}

func readBrandMappings() map[string]string {
	brandMappingsFile, err := ioutil.ReadFile("../brandMappings.json")
	if err != nil {
		log.Errorf("Couldn't read brand mapping configuration: %v\n", err)
		os.Exit(1)
	}
	var brandMappings map[string]string
	err = json.Unmarshal(brandMappingsFile, &brandMappings)
	if err != nil {
		log.Errorf("Couldn't unmarshal brand mapping configuration: %v\n", err)
		os.Exit(1)
	}
	return brandMappings
}

func (eomfile EomFile) initType(txID string) (EomFile, error) {
	contentType := eomfile.ContentType
	contentSrc := eomfile.Source.SourceCode

	if contentSrc == "ContentPlaceholder" && contentType == "EOM::CompoundStory" {
		uuid, err := eomfile.resolveUUID(txID)
		if err != nil {
			return EomFile{}, err
		}

		if uuid != "" {
			eomfile.Type = "EOM::CompoundStory_Internal_CPH"
			eomfile.UUID = uuid
		} else {
			eomfile.Type = "EOM::CompoundStory_External_CPH"
		}

		log.Infof("results [%v] ....", eomfile.Type)
		return eomfile, nil
	}
	eomfile.Type = eomfile.ContentType
	return eomfile, nil
}

func (eomfile EomFile) resolveUUID(txID string) (string, error) {
	attributes, err := buildAttributes(eomfile.Attributes)
	if err != nil {
		return "", err
	}

	uuid := ""
	if isBlogCategory(attributes) {
		resolvedUuid, err := iResolver.ResolveIdentifier(attributes.ServiceId, attributes.RefField, txID)
		if err != nil {
			return "", fmt.Errorf("Couldn't resolve blog uuid %v", err)
		}
		uuid = resolvedUuid
	}
	return uuid, nil
}

func buildAttributes(attributesXML string) (Attributes, error) {
	var attrs Attributes
	if err := xml.Unmarshal([]byte(attributesXML), &attrs); err != nil {
		return Attributes{}, err
	}
	return attrs, nil
}

func isBlogCategory(attributes Attributes) bool {
	for _, c := range blogCategories {
		if c == attributes.Category {
			return true
		}
	}
	return false
}

func (eomfile EomFile) Initialize(binaryContent []byte, txID string) (Content, error) {
	eomfile.BinaryContent = binaryContent
	eomfileInitilized, err := eomfile.initType(txID)
	return Content(eomfileInitilized), err
}

func (eomfile EomFile) Validate(externalValidationEndpoint string, txID string, username string, password string) ValidationResponse {
	validationParam := validationParam{
		eomfile.BinaryContent,
		externalValidationEndpoint,
		username,
		password,
		txID,
		eomfile.GetUUID(),
		eomfile.GetType(),
	}

	return doExternalValidation(
		validationParam,
		eomfile.isValid,
		eomfile.isMarkedDeleted,
	)
}

func (eomfile EomFile) isValid(status int) bool {
	return status == http.StatusOK || status == http.StatusNotFound
}

func (eomfile EomFile) isMarkedDeleted(status ...int) bool {
	return len(status) == 1 && status[0] == http.StatusNotFound
}

func (eomfile EomFile) GetType() string {
	return eomfile.Type
}

func (eomfile EomFile) GetUUID() string {
	return eomfile.UUID
}
