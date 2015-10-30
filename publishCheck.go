package main

import (
	"fmt"
	"net/http"
)

type PublishCheck struct {
	Metric        PublishMetric
	Threshold     int
	CheckInterval int
	ResultSink    chan PublishMetric
}

type EndpointSpecificCheck interface {
	isNewPublish(tid string)
}

type ContentCheck struct{}
type S3Check struct{}

func NewPublishCheck(pm PublishMetric, t int, ci int, rs chan PublishMetric) *PublishCheck {
	return &PublishCheck{pm, t, ci, rs}
}

func (pc PublishCheck) DoCheck() bool {
	info.Printf("Running check for UUID [%v]\n", pc.Metric.UUID)
	resp, err := http.Get(pc.Metric.endpoint.String() + pc.Metric.UUID)
	if err != nil || resp.StatusCode != 200 {
		return false
	} else {
		return true
	}

	//read body
	//extract publishReference, which is the publish TID
	//if we have 200, and publish tid = our tid -> ok
	//if publish tid != our tid -> it;'s a delete or update
	//also get TID here, we don't have it
}

func (c ContentCheck) isNewPublish(tid string) {
	fmt.Println("function isPublishedContent parameter:", tid)
}

func (s S3Check) isNewPublish(tid string) {
	fmt.Println("function isPublishedS3 parameters:", tid)
}

//key is the endpoint alias from the config
var endpointNewPublishChecks = map[string]EndpointSpecificCheck{
	"content": ContentCheck{},
	"S3":      S3Check{},
}
