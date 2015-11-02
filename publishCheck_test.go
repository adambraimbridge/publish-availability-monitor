package main

import (
	"fmt"
	"testing"
)

func TestIsCurrentOperationFinished_S3Check(t *testing.T) {
	s3Check := &S3Check{}
	if !s3Check.isCurrentOperationFinished("tid", []byte("")) {
		t.Errorf("Expected: true. Actual: false")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_InvalidContent(t *testing.T) {
	contentCheck := &ContentCheck{}

	testResponse := `{ "uuid" : "1234-1234"`

	if contentCheck.isCurrentOperationFinished("tid", []byte(testResponse)) {
		t.Errorf("Expected error.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_CurrentOperation(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)
	fmt.Println(testResponse)
	if !contentCheck.isCurrentOperationFinished(currentTid, []byte(testResponse)) {
		t.Error("Expected success.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_NotCurrentOperation(t *testing.T) {
	contentCheck := &ContentCheck{}

	currentTid := "tid_1234"
	testResponse := `{ "uuid" : "1234-1234", "publishReference" : "tid_1235"}`
	fmt.Println(testResponse)
	if contentCheck.isCurrentOperationFinished(currentTid, []byte(testResponse)) {
		t.Error("Expected failure.")
	}
}
