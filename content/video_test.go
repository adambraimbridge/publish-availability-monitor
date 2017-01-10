package content

import (
	"testing"
)

var videoValid = Video{
	UUID:        "e28b12f7-9796-3331-b030-05082f0b8157",
	Id:          "4966650664001",
	Name:        "the-dark-knight.mp4",
	UpdatedAt:   "2016-06-01T21:40:19.120Z",
	PublishedAt: "2016-06-01T21:40:19.120Z",
}

func TestIsVideoValid_Valid(t *testing.T) {
	if !videoValid.IsValid("", "", "", "") {
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
