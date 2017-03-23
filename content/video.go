package content

import (
	"regexp"
)

const videoType = "video"

var idRegexp, _ = regexp.Compile("^\\d+$")

type Video struct {
	UUID          string `json:"uuid"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	UpdatedAt     string `json:"updated_at"`
	PublishedAt   string `json:"published_at"`
	BinaryContent []byte `json:"-"` //This field is for internal application usage
}

func (video Video) Initialize(binaryContent []byte) Content {
	video.BinaryContent = binaryContent
	return video
}

func (v Video) IsValid(externalValidationEndpoint string, txId string, username string, password string) bool {
	contentUUID := v.UUID

	if !isUUIDValid(contentUUID) {
		warnLogger.Printf("Video invalid: invalid UUID: [%s]", contentUUID)
		return false
	}

	if !idRegexp.MatchString(v.Id) {
		warnLogger.Printf("Video invalid: invalid ID: [%s]", v.Id)
		return false
	}

	return isExternalValidationSuccessful(v.BinaryContent, externalValidationEndpoint, username, password, txId, contentUUID, v.GetType())
}

func (v Video) IsMarkedDeleted() bool {
	if v.PublishedAt != "" || v.UpdatedAt != "" {
		return false
	}
	return true
}

func (v Video) GetType() string {
	return videoType
}

func (v Video) GetUUID() string {
	return v.UUID
}
