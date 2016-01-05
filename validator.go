package main

import (
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

const syntheticPrefix = "SYNTHETIC"

func isMessagePastPublishSLA(date time.Time, threshold int) bool {
	passedSLA := date.Add(time.Duration(threshold) * time.Second)
	return time.Now().After(passedSLA)
}

func isSyntheticMessage(tid string) bool {
	return strings.HasPrefix(tid, syntheticPrefix)
}

func isUUIDValid(contentUUID string) bool {
	parsedUUID, err := uuid.FromString(contentUUID)
	if err != nil {
		warn.Printf("Cannot parse UUID [%v], error: [%v]", contentUUID, err.Error())
		return false
	}
	return contentUUID == parsedUUID.String()
}
