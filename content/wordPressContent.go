package content

import "net/http"

const notFoundError = "Not found."

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

func (wordPressMessage WordPressMessage) IsValid(extValEndpoint string, txId string, username string, password string) bool {
	contentUUID := wordPressMessage.Post.UUID
	if !isUUIDValid(contentUUID) {
		warnLogger.Printf("WordPress message invalid: invalid UUID: [%s]", contentUUID)
		return false
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
		wordpressStatusCheck,
	)
}

func (wordPressMessage WordPressMessage) IsMarkedDeleted() bool {
	return wordPressMessage.Status == "error" && wordPressMessage.Error == notFoundError
}

func (wordPressMessage WordPressMessage) GetType() string {
	return wordPressMessage.Post.Type
}

func (wordPressMessage WordPressMessage) GetUUID() string {
	return wordPressMessage.Post.UUID
}

func wordpressStatusCheck(status int) bool {
	return status != http.StatusUnprocessableEntity
}
