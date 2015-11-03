package main

import (
	"log"
	"os"
)

// SplunkFeeder implements MetricDestination interface to send PublishMetrics to Splunk.
// This is achieved by writing the metric into a file which is indexed by Splunk.
type SplunkFeeder struct {
	MetricLog *log.Logger
}

// NewSplunkFeeder returns a SplunkFeeder which will write the PublishMetrics to the file at filePath.
// If the file exists, it will be appended to.
func NewSplunkFeeder(filePath string) *SplunkFeeder {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Cannot open/create file [%v] : [%v]", filePath, err.Error())
		return nil
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	return &SplunkFeeder{logger}

}

// Send logs pm into a file.
func (sf SplunkFeeder) Send(pm PublishMetric) {
	sf.MetricLog.Printf("UUID=%v publishDate=%v publishOk=%v duration=%v endpoint=%v ",
		pm.UUID, pm.publishDate.UnixNano(), pm.publishOK, pm.publishInterval.upperBound, pm.config.Alias)
}
