package main

var validWordPressTypes []string

func init() {
	validWordPressTypes = []string{
		"post",
		"webchat-markets-live",
		"webchat-live-blogs",
		"webchat-live-qa",
	}
}

// WordPressMessage models messages from Wordpress
type WordPressMessage struct {
	Status      string `json:"status"`
	APIURL      string `json:"apiUrl"`
	Error       string `json:"error"`
	Post        *Post  `json:"post"`
	PreviousURL string `json:"previousUrl"`
}

// Post models WordPress content
// neglect unused fields (e.g. slug, title, content, etc)
type Post struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

func (wordPressMessage WordPressMessage) isValid() bool {
	contentUUID := wordPressMessage.Post.UUID
	if !isUUIDValid(contentUUID) {
		warn.Printf("WordPress message invalid: invalid UUID: [%s]", contentUUID)
		return false
	}

	apiURL := wordPressMessage.APIURL
	if !isValidBrand(apiURL) {
		warn.Printf("WordPress message invalid: failed to resolve brand for uri [%s].", apiURL)
		return false
	}

	contentType := wordPressMessage.Post.Type
	for _, validType := range validWordPressTypes {
		if contentType == validType {
			return true
		}
	}
	warn.Printf("WordPress message invalid: unexpected content type: [%s]", contentType)
	return false
}

func (wordPressMessage WordPressMessage) isMarkedDeleted() bool {
	if wordPressMessage.Post == nil {
		return true
	}
	return false
}

func (wordPressMessage WordPressMessage) getType() string {
	return wordPressMessage.Post.Type
}

func (wordPressMessage WordPressMessage) getUUID() string {
	return wordPressMessage.Post.UUID
}
