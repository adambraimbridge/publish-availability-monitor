package checks

import (
	"strings"
)

func IsIgnorableMessage(tid string) bool {
	return strings.HasPrefix(tid, "SYNTHETIC")
}
