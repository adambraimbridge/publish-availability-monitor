package main

import (
	"testing"
)

func TestGenerateImageSetUUID(t *testing.T) {
	expectedImagesetUUIDString := "e440b756-16b9-4520-125b-772aa0ca37ea"

	imageUUIDString := "e440b756-16b9-4520-8c3d-e0b07abe9ca3"
	imageUUID, _ := NewUUIDFromString(imageUUIDString)

	imagesetUUID, err := GenerateImageSetUUID(*imageUUID)

	if err != nil {
		t.Error("Returned error for valid input")
	}
	actualImagesetUUIDString := imagesetUUID.String()

	if expectedImagesetUUIDString != actualImagesetUUIDString {
		t.Error("Imageset UUID was not generated correctly.")
	}
}
