package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type PublishCheck struct {
	Metric        PublishMetric
	Threshold     int
	CheckInterval int
	ResultSink    chan PublishMetric
	tid           string
}

type EndpointSpecificCheck interface {
	isCurrentOperationFinished(ourTid string, resp []byte) bool
}

type ContentCheck struct{}
type S3Check struct{}

func NewPublishCheck(pm PublishMetric, t int, ci int, rs chan PublishMetric, tid string) *PublishCheck {
	return &PublishCheck{pm, t, ci, rs, tid}
}

func (pc PublishCheck) DoCheck() bool {
	info.Printf("Running check for UUID [%v]\n", pc.Metric.UUID)
	resp, err := http.Get(pc.Metric.endpoint.String() + pc.Metric.UUID)

	if err != nil || resp.StatusCode != 200 {
		return false
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		warn.Printf("Cannot read response: [%s]", err.Error())
		return false
	}

	check := endpointSpecificChecks[pc.Metric.config.Alias]
	if check == nil {
		warn.Printf("No check for endpoint %s.", pc.Metric.config.Alias)
		return false
	}
	return check.isCurrentOperationFinished(pc.tid, data)
}

func (c ContentCheck) isCurrentOperationFinished(tid string, resp []byte) bool {
	info.Println("Content isCurrentOperationFinished() check")
	var jsonResp map[string]interface{}

	err := json.Unmarshal(resp, &jsonResp)
	if err != nil {
		warn.Printf("Cannot unmarshal JSON response: [%s]", err.Error())
		return false
	}

	return jsonResp["publishReference"] == tid
}

func (s S3Check) isCurrentOperationFinished(tid string, resp []byte) bool {
	info.Println("S3 isCurrentOperationFinished() check")
	return true
}

//key is the endpoint alias from the config
var endpointSpecificChecks = map[string]EndpointSpecificCheck{
	"content": ContentCheck{},
	"S3":      S3Check{},
}
