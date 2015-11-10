package main

import (
	"testing"
)

func TestValidType(testing *testing.T) {
	var tests = []struct {
		validTypes []string
		eomType    string
		expected   bool
	}{
		{
			[]string{"Image", "EOM:WebContainer"},
			"EOM:CompoundStory",
			false,
		},
		{
			[]string{"Image", "EOM:WebContainer"},
			"EOM:WebContainer",
			true,
		},
	}

	for _, t := range tests {
		actual := validType(t.validTypes, t.eomType)
		if actual != t.expected {
			testing.Errorf("Test Case: %v\nActual: %v", t, actual)
		}
	}
}
