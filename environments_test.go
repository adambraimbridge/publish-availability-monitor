package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEtcdValues(t *testing.T) {
	environments := make(map[string]Environment)
	parseEnvironmentsIntoMap("t1:https://t1.example.org,t2:https://t2.example.com", "t1:user1:pass1,t2:user2:pass2", environments)

	t1 := environments["t1"]
	assert.Equal(t, "t1", t1.Name, "environment name")
	assert.Equal(t, "https://t1.example.org", t1.ReadUrl, "environment read url")
	assert.Equal(t, "user1", t1.Username, "environment username")
	assert.Equal(t, "pass1", t1.Password, "environment password")

	t2 := environments["t2"]
	assert.Equal(t, "t2", t2.Name, "environment name")
	assert.Equal(t, "https://t2.example.com", t2.ReadUrl, "environment read url")
	assert.Equal(t, "user2", t2.Username, "environment username")
	assert.Equal(t, "pass2", t2.Password, "environment password")

	assert.Equal(t, len(environments), 2, "environments")
}

func TestParseEtcdUnauthValues(t *testing.T) {
	environments := make(map[string]Environment)
	parseEnvironmentsIntoMap("t1:https://t1.example.org,t2:https://t2.example.com", "t2:user2:pass2", environments)

	t1 := environments["t1"]
	assert.Equal(t, "t1", t1.Name, "environment name")
	assert.Equal(t, "https://t1.example.org", t1.ReadUrl, "environment read url")
	assert.Equal(t, "", t1.Username, "environment username")
	assert.Equal(t, "", t1.Password, "environment password")

	t2 := environments["t2"]
	assert.Equal(t, "t2", t2.Name, "environment name")
	assert.Equal(t, "https://t2.example.com", t2.ReadUrl, "environment read url")
	assert.Equal(t, "user2", t2.Username, "environment username")
	assert.Equal(t, "pass2", t2.Password, "environment password")

	assert.Equal(t, len(environments), 2, "environments")
}

func TestParseEmptyEtcdValues(t *testing.T) {
	environments := make(map[string]Environment)
	parseEnvironmentsIntoMap("", "", environments)

	assert.Empty(t, environments, "expected an empty map")
}
