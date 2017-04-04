package content

import (
	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var wordpressContentMarkedDeletedTrue = WordPressMessage{
	Status: "error", Error: "Not found.",
}

var wordpressContentMarkedDeletedFalse = WordPressMessage{
	Status: "ok", Post: Post{},
}

func TestIsMarkedDeleted_True(t *testing.T) {
	valRes := wordpressContentMarkedDeletedTrue.Validate("", "", "", "")
	if !valRes.IsMarkedDeleted {
		t.Error("Expected True, the story IS marked deleted")
	}
}

func TestIsMarkedDeleted_False(t *testing.T) {
	valRes := wordpressContentMarkedDeletedFalse.Validate("", "", "", "")
	if valRes.IsMarkedDeleted {
		t.Error("Expected False, the wordPress article IS NOT marked deleted")
	}
}

func TestIsValid_FailedExternalValidation422(t *testing.T) {
	var wordpressMessage = WordPressMessage{
		Status:        "ok",
		Error:         "",
		Post:          Post{UUID: "e28b12f7-9796-3331-b030-05082f0b8157"},
		PreviousURL:   "",
		BinaryContent: []byte{},
	}

	txId := "tid_1234"
	pamTxId := checks.ConstructPamTxId(txId)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxId, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, wordpressMessage.BinaryContent, reqBody)

		w.WriteHeader(http.StatusUnprocessableEntity)
	}))

	if wordpressMessage.IsValid(testServer.URL+"/map", txId, "", "") {
		t.Error("Wordpress should fail external validation.")
	}
}
