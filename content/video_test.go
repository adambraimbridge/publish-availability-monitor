package content

import (
	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var videoValid = Video{
	UUID:        "e28b12f7-9796-3331-b030-05082f0b8157",
	Id:          "4966650664001",
	Name:        "the-dark-knight.mp4",
	UpdatedAt:   "2016-06-01T21:40:19.120Z",
	PublishedAt: "2016-06-01T21:40:19.120Z",
	binaryContent: []byte(
		`{
		"id":"4922311929001",
		"uuid":"2a304b92-7d99-34bd-ad7b-2d781bfcedb8",
		"name":"the-dark-knight.mp4",
		"UpdatedAt":"2016-06-01T21:40:19.120Z",
		"PublishedAt":"2016-06-01T21:40:19.120Z"
	}`),
}

func TestIsVideoValid_Valid(t *testing.T) {
	txId := "tid_1234"
	pamTxId := checks.ConstructPamTxId(txId)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/map", req.RequestURI)
		assert.Equal(t, pamTxId, req.Header.Get("X-Request-Id"))

		defer req.Body.Close()
		reqBody, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, videoValid.binaryContent, reqBody)
	}))

	if !videoValid.IsValid(testServer.URL+"/map", txId, "", "") {
		t.Error("Video should be valid.")
	}
}

var videoNoId = Video{
	UUID:        "e28b12f7-9796-3331-b030-05082f0b8157",
	Name:        "the-dark-knight.mp4",
	UpdatedAt:   "2016-06-01T21:40:19.120Z",
	PublishedAt: "2016-06-01T21:40:19.120Z",
}

func TestIsVideoValid_NoId(t *testing.T) {
	if videoNoId.IsValid("", "", "", "") {
		t.Error("Video should be invalid as it has no Id.")
	}
}

var videoNoUUID = Video{
	Id:          "4966650664001",
	Name:        "the-dark-knight.mp4",
	UpdatedAt:   "2016-06-01T21:40:19.120Z",
	PublishedAt: "2016-06-01T21:40:19.120Z",
}

func TestIsVideoValid_NoUUID(t *testing.T) {
	if videoNoUUID.IsValid("", "", "", "") {
		t.Error("Video should be invalid as it has no uuid.")
	}
}

var videoNoDates = Video{
	UUID: "e28b12f7-9796-3331-b030-05082f0b8157",
	Id:   "4966650664001",
	Name: "the-dark-knight.mp4",
}

func TestIsDeleted_NoDates(t *testing.T) {
	if !videoNoDates.IsMarkedDeleted() {
		t.Error("Video should be evaluated as deleted as it has no dates in it.")
	}
}

var videoOneDateOnly = Video{
	UUID:      "e28b12f7-9796-3331-b030-05082f0b8157",
	Id:        "4966650664001",
	Name:      "the-dark-knight.mp4",
	UpdatedAt: "2016-06-01T21:40:19.120Z",
}

func TestIsDeleted_OneDateOnly(t *testing.T) {
	if videoOneDateOnly.IsMarkedDeleted() {
		t.Error("Video should be evaluated as published as it has one date in it.")
	}
}
