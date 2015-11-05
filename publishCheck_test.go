package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestIsCurrentOperationFinished_S3Check_Finished(t *testing.T) {
	s3Check := &S3Check{}
	if !s3Check.isCurrentOperationFinished(PublishCheck{}, &http.Response{StatusCode: 200}) {
		t.Errorf("Expected: true. Actual: false")
	}
}

func TestIsCurrentOperationFinished_S3Check_NotFinished(t *testing.T) {
	s3Check := &S3Check{}
	if s3Check.isCurrentOperationFinished(PublishCheck{}, &http.Response{StatusCode: 404}) {
		t.Errorf("Expected: false. Actual: True")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_InvalidContent(t *testing.T) {
	contentCheck := &ContentCheck{}
	testResponse := `{ "uuid" : "1234-1234"`

	if contentCheck.isCurrentOperationFinished(buildPublishCheck(false, "tid"), buildResponse(200, testResponse)) {
		t.Errorf("Expected error.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_CurrentOperation(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)

	if !contentCheck.isCurrentOperationFinished(buildPublishCheck(false, currentTid), buildResponse(200, testResponse)) {
		t.Error("Expected success.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_NotCurrentOperation(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := `{ "uuid" : "1234-1234", "publishReference" : "tid_1235"}`
	fmt.Println(testResponse)
	if contentCheck.isCurrentOperationFinished(buildPublishCheck(false, currentTid), buildResponse(200, testResponse)) {
		t.Error("Expected failure.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_MarkedDeleted_Finished(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)

	if !contentCheck.isCurrentOperationFinished(buildPublishCheck(true, currentTid), buildResponse(404, testResponse)) {
		t.Error("Expected success.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_MarkedDeleted_NotFinished(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)

	if contentCheck.isCurrentOperationFinished(buildPublishCheck(true, currentTid), buildResponse(200, testResponse)) {
		t.Error("Expected failure.")
	}
}

func buildPublishCheck(isMarkedDeleted bool, tid string) PublishCheck {
	return PublishCheck{
		Metric: PublishMetric{
			isMarkedDeleted: isMarkedDeleted,
			tid:             tid,
		},
	}
}

func buildResponse(statusCode int, content string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       nopCloser{bytes.NewBuffer([]byte(content))},
	}
}

// this is necessary to be able to build an http.Response
// the body has to be a ReadCloser
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
