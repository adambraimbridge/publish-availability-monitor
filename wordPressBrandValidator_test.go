package main

import (
	"testing"
)

func TestIsValidBrand_Valid(t *testing.T) {
	if !isValidBrand("http://blogs.ft.com/tech-blog/post?id=123456") {
		t.Errorf("Expected True.")
	}
}

func TestIsValidBrand_NotValid(t *testing.T) {
	if isValidBrand("http://blogs.ft.com/foobar/post?id=123456") {
		t.Errorf("Expected False.")
	}
}
