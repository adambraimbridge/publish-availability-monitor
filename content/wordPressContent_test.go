package content

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/stretchr/testify/assert"
)

func TestIsValid_ExternalValidationInvalidArticle422(t *testing.T) {
	var wordpressMessage = WordPressMessage{
		Status:        "ok",
		Error:         "",
		Post:          Post{UUID: "e28b12f7-9796-3331-b030-05082f0b8157"},
		PreviousURL:   "",
		BinaryContent: []byte{},
	}

	txID := "tid_1234"
	pamTxID := checks.ConstructPamTxID(txID)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxID, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, wordpressMessage.BinaryContent, reqBody)

		w.WriteHeader(http.StatusUnprocessableEntity)
	}))

	valResp := wordpressMessage.Validate(testServer.URL+"/map", txID, "", "")
	if valResp.IsValid {
		t.Error("Wordpress should fail external validation.")
	}
}

func TestIsValid_ExternalValidationMarkedAsDeleted404(t *testing.T) {
	var wordpressMessage = WordPressMessage{
		Status:        "ok",
		Error:         "",
		Post:          Post{UUID: "e28b12f7-9796-3331-b030-05082f0b8157"},
		PreviousURL:   "",
		BinaryContent: []byte{},
	}

	txID := "tid_1234"
	pamTxID := checks.ConstructPamTxID(txID)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxID, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, wordpressMessage.BinaryContent, reqBody)

		w.WriteHeader(http.StatusNotFound)
	}))

	valResp := wordpressMessage.Validate(testServer.URL+"/map", txID, "", "")
	if !valResp.IsValid {
		t.Error("Wordpress article marked as deleted shouldn't fail external validation.")
	}

	if !valResp.IsMarkedDeleted {
		t.Error("Wordpress article should be marked as deleted.")
	}
}
