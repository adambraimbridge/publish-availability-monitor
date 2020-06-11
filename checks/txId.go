package checks

import (
	"strings"
)

func ConstructPamTxID(txID string) string {
	if strings.HasPrefix(txID, "tid_") {
		txID = txID[:4] + "pam_" + txID[4:]
	}

	return txID
}
