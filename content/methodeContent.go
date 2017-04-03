package content

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/xml"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	log "github.com/Sirupsen/logrus"
	xmlpath "gopkg.in/xmlpath.v1"
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
		log.Infof("results [%v] ....", eomfile.Type)
		return eomfile
	}
	eomfile.Type = eomfile.ContentType
	return eomfile
}

func (eomfile EomFile) IsValid(externalValidationEndpoint string, txID string, username string, password string) bool {
	contentUUID := eomfile.UUID
	if !isUUIDValid(contentUUID) {
		log.Warnf("Eomfile invalid: invalid UUID: [%s]. transaction_id=[%s]", contentUUID, txID)
		return false
	}

	return isExternalValidationSuccessful(eomfile, externalValidationEndpoint, txID, username, password)
}

func (eomfile EomFile) IsMarkedDeleted() bool {
	if eomfile.Type == "Image" || eomfile.Type == "EOM::WebContainer" {
		return false
	}
	markDeletedFlag, ok := GetXPathValue(eomfile.Attributes, eomfile, markDeletedFlagXPath)
	if !ok {
		log.Warnf("Eomfile with uuid=[%s]: Cannot match node in XML using xpath [%v]", eomfile.UUID, markDeletedFlagXPath)
		return false
	}
	log.Infof("Eomfile with uuid=[%s]: MarkAsDeletedFlag: [%v]", eomfile.UUID, markDeletedFlag)
	return markDeletedFlag == "True"
}

func (eomfile EomFile) GetType() string {
	return eomfile.Type
}

func (eomfile EomFile) GetUUID() string {
	return eomfile.UUID
}

func GetXPathValue(xml string, eomfile EomFile, lookupPath string) (string, bool) {
	path := xmlpath.MustCompile(lookupPath)
	root, err := xmlpath.Parse(strings.NewReader(xml))
	if err != nil {
		log.Warnf("Cannot parse XML of eomfile with uuid=[%s] using xpath [%v], error: [%v]", eomfile.UUID, lookupPath, err.Error())
		return "", false
	}
	xpathValue, ok := path.String(root)
	return xpathValue, ok

}

func isExternalValidationSuccessful(eomfile EomFile, validationURL string, txID, username string, password string) bool {
	if validationURL == "" {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s]. Validation endpoint URL is missing for content type=[%s]", eomfile.UUID, txID, eomfile.Type)
		return false
	}
	marshalled, err := json.Marshal(eomfile)
	if err != nil {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]. Skipping external validation.", eomfile.UUID, txID, err)
		return true
	}

	resp, err := httpCaller.DoCallWithEntity(
		"POST", validationURL,
		username, password,
		checks.ConstructPamTxId(txID),
		"application/json", bytes.NewReader(marshalled))

	if err != nil {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]. Skipping external validation.", eomfile.UUID, txID, err)
		return true
	}
	defer cleanupResp(resp)

	log.Infof("External validation for content uuid=[%s] transaction_id=[%s] received statusCode [%d]", eomfile.UUID, txID, resp.StatusCode)

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]", eomfile.UUID, txID, err)
	}
	if resp.StatusCode != 200 {
		log.Infof("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]", eomfile.UUID, txID, string(bs))
	}
	if resp.StatusCode == 418 {
		return false
	}
	//invalid  contentplaceholder (link file) will not be published so do not monitor
	if resp.StatusCode == 422 {
		return false
	}

	return true
}

func cleanupResp(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		log.Warnf("[%v]", err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Warnf("[%v]", err)
	}
}
