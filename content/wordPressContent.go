package content

import "net/http"

// WordPressMessage models messages from Wordpress
type WordPressMessage struct {
	Status        string `json:"status"`
	Error         string `json:"error"`
	Post          Post   `json:"post"`
	PreviousURL   string `json:"previousUrl"`
	BinaryContent []byte `json:"-"` //This field is for internal application usage
}

// Post models WordPress content
// neglect unused fields (e.g. id, slug, title, content, etc)
type Post struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
	Url  string `json:"url"`
}

func (wordPressMessage WordPressMessage) Initialize(binaryContent []byte) Content {
	wordPressMessage.BinaryContent = binaryContent
	return wordPressMessage
}

func (wordPressMessage WordPressMessage) Validate(extValEndpoint string, txId string, username string, password string) ValidationResponse {
	contentUUID := wordPressMessage.Post.UUID
	if !isUUIDValid(contentUUID) {
		warnLogger.Printf("WordPress message invalid: invalid UUID: [%s]", contentUUID)
		return ValidationResponse{IsValid: false, IsMarkedDeleted: wordPressMessage.isMarkedDeleted(0)}
	}

	validationParam := validationParam{
		wordPressMessage.BinaryContent,
		extValEndpoint,
		username,
		password,
		txId,
		wordPressMessage.GetUUID(),
		wordPressMessage.GetType(),
	}

	return doExternalValidation(
		validationParam,
		wordPressMessage.isValid,
		wordPressMessage.isMarkedDeleted,
	)
}

func (wordPressMessage WordPressMessage) isValid(status int) bool {
	return status != http.StatusUnprocessableEntity
}

func (wordPressMessage WordPressMessage) isMarkedDeleted(status ...int) bool {
	return len(status) == 1 && status[0] == http.StatusNotFound
}

func (wordPressMessage WordPressMessage) GetType() string {
	return wordPressMessage.Post.Type
}

func (wordPressMessage WordPressMessage) GetUUID() string {
	return wordPressMessage.Post.UUID
}
