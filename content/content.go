package content

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/publish-availability-monitor/checks"
	log "github.com/Sirupsen/logrus"
)

// Content is the interface for different type of contents from different CMSs.
type Content interface {
	Initialize(binaryContent []byte) Content
	Validate(externalValidationEndpoint string, txID string, username string, password string) ValidationResponse
	GetType() string
	GetUUID() string
}

type ValidationResponse struct {
	IsValid         bool
	IsMarkedDeleted bool
}

type validationParam struct {
	binaryContent []byte
	validationURL string
	username      string
	password      string
	txID          string
	uuid          string
	contentType   string
}

const systemIDKey = "Origin-System-Id"

var httpCaller checks.HttpCaller

func init() {
	httpCaller = checks.NewHttpCaller(10)
}

// UnmarshalContent unmarshals the message body into the appropriate content type based on the systemID header.
func UnmarshalContent(msg consumer.Message) (Content, error) {
	binaryContent := []byte(msg.Body)

	headers := msg.Headers
	systemID := headers[systemIDKey]
	switch systemID {
	case "http://cmdb.ft.com/systems/methode-web-pub":
		var eomFile EomFile

		err := json.Unmarshal(binaryContent, &eomFile)
		if err != nil {
			return nil, err
		}
		xml.Unmarshal([]byte(eomFile.Attributes), &eomFile.Source)

		return eomFile.Initialize(binaryContent), err
	case "http://cmdb.ft.com/systems/wordpress":
		var wordPressMsg WordPressMessage
		err := json.Unmarshal(binaryContent, &wordPressMsg)
		return wordPressMsg.Initialize(binaryContent), err
	case "http://cmdb.ft.com/systems/next-video-editor":
		var video Video
		err := json.Unmarshal(binaryContent, &video)
		return video.Initialize(binaryContent), err
	default:
		return nil, fmt.Errorf("Unsupported content with system ID: [%s].", systemID)
	}
}

func doExternalValidation(p validationParam, validCheck func(int) bool, deletedCheck func(...int) bool) ValidationResponse {
	if p.validationURL == "" {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s]. Validation endpoint URL is missing for content type=[%s]", p.uuid, p.txID, p.contentType)
		return ValidationResponse{false, deletedCheck()}
	}

	resp, err := httpCaller.DoCall(checks.Config{
		HttpMethod: "POST", Url: p.validationURL, Username: p.username, Password: p.password,
		TxId:        checks.ConstructPamTxId(p.txID),
		ContentType: "application/json", Entity: bytes.NewReader(p.binaryContent)})

	if err != nil {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s] validationURL=[%s], creating validation request error: [%v]. Skipping external validation.", p.uuid, p.txID, p.validationURL, err)
		return ValidationResponse{true, deletedCheck()}
	}
	defer cleanupResp(resp)

	log.Infof("External validation for content uuid=[%s] transaction_id=[%s] validationURL=[%s], received statusCode [%d]", p.uuid, p.txID, p.validationURL, resp.StatusCode)

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("External validation for content uuid=[%s] transaction_id=[%s] validationURL=[%s], reading response body error: [%v]", p.uuid, p.txID, p.validationURL, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		log.Infof("External validation for content uuid=[%s] transaction_id=[%s] validationURL=[%s], received statusCode [%d], received error: [%v]", p.uuid, p.txID, p.validationURL, resp.StatusCode, string(bs))
	}

	if resp.StatusCode == http.StatusNotFound {
		log.Infof("External validation for content uuid=[%s] transaction_id=[%s] validationURL=[%s], received statusCode [%d], content is marked as deleted.", p.uuid, p.txID, p.validationURL, resp.StatusCode)
	}

	return ValidationResponse{validCheck(resp.StatusCode), deletedCheck(resp.StatusCode)}
}

func cleanupResp(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		log.Warnf("External validation cleanup failed with error: [%v]", err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Warnf("External validation cleanup failed with error: [%v]", err)
	}
}
