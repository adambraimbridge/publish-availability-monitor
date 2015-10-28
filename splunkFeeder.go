package main

import (
	"log"
	"os"
)

type SplunkFeeder struct {
	MetricLog *log.Logger
}

func NewSplunkFeeder(filePath string) *SplunkFeeder {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Cannot open/create file [%v] : [%v]", filePath, err.Error())
		return nil
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LUTC)
	return &SplunkFeeder{logger}

}
func (sf SplunkFeeder) Send(pm PublishMetric) {
	sf.MetricLog.Printf("UUID=%v publishDate=%v publishOk=%v duration=%v endpoint=%v ",
		pm.UUID, pm.publishDate.Unix(), pm.publishOK, pm.publishInterval.upperBound, pm.config.Alias)
}
