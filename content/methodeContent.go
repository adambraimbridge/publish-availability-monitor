package content

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"launchpad.net/xmlpath"
)

const SourceXPath = "//ObjectMetadata/EditorialNotes/Sources/Source/SourceCode"
const markDeletedFlagXPath = "//ObjectMetadata/OutputChannels/DIFTcom/DIFTcomMarkDeleted"

// EomFile models Methode content
type EomFile struct {
	Attributes       string        `json:"attributes"`
	LinkedObjects    []interface{} `json:"linkedObjects"`
	SystemAttributes string        `json:"systemAttributes"`
	Type             string        `json:"type"`
	UsageTickets     string        `json:"usageTickets"`
	UUID             string        `json:"uuid"`
	Value            string        `json:"value"`
	WorkflowStatus   string        `json:"workflowStatus"`
}

func (eomfile EomFile) IsValid(externalValidationEndpoint string, txID string, username string, password string) bool {
	contentUUID := eomfile.UUID
	if !isUUIDValid(contentUUID) {
		warnLogger.Printf("Eomfile invalid: invalid UUID: [%s]. transaction_id=[%s]", contentUUID, txID)
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

func GetXPathValue(xml string, eomfile EomFile, lookupPath string) (string, bool) {
	path := xmlpath.MustCompile(lookupPath)
	root, err := xmlpath.Parse(strings.NewReader(xml))
	if err != nil {
		warnLogger.Printf("Cannot parse XML of eomfile with uuid=[%s] using xpath [%v], error: [%v]", eomfile.UUID, lookupPath, err.Error())
		return "", false
	}
	xpathValue, ok := path.String(root)
	return xpathValue, ok

}

func isExternalValidationSuccessful(eomfile EomFile, validationURL string, txID, username string, password string) bool {
	if validationURL == "" {
		warnLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s]. Validation endpoint URL is missing for content type=[%s]", eomfile.UUID, txID, eomfile.Type)
		return false
	}
	marshalled, err := json.Marshal(eomfile)
	if err != nil {
		warnLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]. Skipping external validation.", eomfile.UUID, txID, err)
		return true
	}

	resp, err := httpCaller.DoCallWithEntity(
		"POST", validationURL,
		username, password,
		checks.ConstructPamTxId(txID),
		"application/json", bytes.NewReader(marshalled))

	if err != nil {
		warnLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]. Skipping external validation.", eomfile.UUID, txID, err)
		return true
	}
	defer cleanupResp(resp)

	infoLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s] received statusCode [%d]", eomfile.UUID, txID, resp.StatusCode)

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		warnLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]", eomfile.UUID, txID, err)
	}
	if resp.StatusCode != 200 {
		infoLogger.Printf("External validation for content uuid=[%s] transaction_id=[%s] error: [%v]", eomfile.UUID, txID, string(bs))
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
		warnLogger.Printf("[%v]", err)
	}
	err = resp.Body.Close()
	if err != nil {
		warnLogger.Printf("[%v]", err)
	}
}
