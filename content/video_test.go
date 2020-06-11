package content

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/stretchr/testify/assert"
)

func TestIsVideoValid_Valid(t *testing.T) {
	var videoValid = Video{
		ID:            "e28b12f7-9796-3331-b030-05082f0b8157",
		BinaryContent: []byte("valid-json"),
	}

	txID := "tid_1234"
	pamTxID := checks.ConstructPamTxID(txID)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxID, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, videoValid.BinaryContent, reqBody)
	}))

	validationResponse := videoValid.Validate(testServer.URL+"/map", txID, "", "")
	assert.True(t, validationResponse.IsValid, "Video should be valid.")
}

func TestIsVideoValid_NoId(t *testing.T) {
	var videoNoId = Video{}

	validationResponse := videoNoId.Validate("", "", "", "")
	assert.False(t, validationResponse.IsValid, "Video should be invalid as it has no Id.")
}

func TestIsVideoValid_failedExternalValidation(t *testing.T) {
	var videoInvalid = Video{
		ID:            "e28b12f7-9796-3331-b030-05082f0b8157",
		BinaryContent: []byte("invalid-json"),
	}

	txID := "tid_1234"
	pamTxID := checks.ConstructPamTxID(txID)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxID, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, videoInvalid.BinaryContent, reqBody)

		w.WriteHeader(http.StatusBadRequest)
	}))

	validationResponse := videoInvalid.Validate(testServer.URL+"/map", txID, "", "")
	assert.False(t, validationResponse.IsMarkedDeleted, "Video should fail external validation.")
}

func TestIsDeleted(t *testing.T) {
	var videoNoDates = Video{
		ID:      "e28b12f7-9796-3331-b030-05082f0b8157",
		Deleted: true,
	}

	validationResponse := videoNoDates.Validate("", "", "", "")
	assert.True(t, validationResponse.IsMarkedDeleted, "Video should be evaluated as deleted.")
}
