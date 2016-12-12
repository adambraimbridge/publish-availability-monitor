package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Financial-Times/publish-availability-monitor/history"
)

func loadHistory(w http.ResponseWriter, r *http.Request) {
	metricContainer.RLock()
	defer metricContainer.RUnlock()

	if strings.Contains(r.Header.Get("Accept"), "text/plain") {
		writePlainTextHistory(w)
	} else {
		w.Header().Add("Content-Type", "text/html")
		history.WriteHistory(w)
	}
}

func writePlainTextHistory(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/plain")
	for i := len(metricContainer.publishMetrics) - 1; i >= 0; i-- {
		fmt.Fprintf(w, "%d. %v\n\n", len(metricContainer.publishMetrics)-i, metricContainer.publishMetrics[i])
	}
}

func forget(w http.ResponseWriter, r *http.Request) {
	tid := r.FormValue("tid")
	history.Forget(tid)
	loadHistory(w, r)
}
