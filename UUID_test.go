package main

import (
	"testing"
)

func TestNewNameUUIDFromBytes(t *testing.T) {
	expectedUUID := "2b588635-d83d-3d6f-9e66-979ada74ab49"
	uuid := NewNameUUIDFromBytes([]byte("imageset"))

	actualUUID := uuid.String()

	if actualUUID != expectedUUID {
		t.Error("UUID not created correctly")
	}
}

func TestNewUUIDFromString_validString(t *testing.T) {
	expectedUUID := "2b588635-d83d-3d6f-9e66-979ada74ab49"
	uuid, err := NewUUIDFromString(expectedUUID)

	if err != nil {
		t.Error("Error returned for valid input")
	}

	actualUUID := uuid.String()

	if actualUUID != expectedUUID {
		t.Error("UUID not created correctly")
	}
}

func TestNewUUIDFromString_invalidString(t *testing.T) {
	expectedUUID := "not-a-valid-uuid-d83d-3d6f-9e66-979ada74ab49"
	_, err := NewUUIDFromString(expectedUUID)

	if err == nil {
		t.Error("Error expected for invalid input")
	}

}
