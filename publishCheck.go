package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// PublishCheck performs an availability  check on a piece of content, at a
// given endpoint, and returns whether the check was successful or not.
// Holds all the information necessary to check content availability
// at an endpoint, as well as store and send the results of the check.
type PublishCheck struct {
	Metric        PublishMetric
	Threshold     int
	CheckInterval int
	ResultSink    chan PublishMetric
}

// EndpointSpecificCheck is the interface which defines a method which determines
// the state of the operation we are currently checking.
type EndpointSpecificCheck interface {
	isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool
}

// ContentCheck implements the EndpointSpecificCheck interface to check operation
// status for the content endpoint.
type ContentCheck struct{}

// S3Check implements the EndpointSpecificCheck interface to check operation
// status for the S3 endpoint.
type S3Check struct{}

// NewPublishCheck returns a PublishCheck ready to perform a check for pm.UUID, at the
// pm.endpoint.
func NewPublishCheck(pm PublishMetric, t int, ci int, rs chan PublishMetric) *PublishCheck {
	return &PublishCheck{pm, t, ci, rs}
}

// DoCheck performs an availability check on a piece of content at a certain
// endpoint, applying endoint-specific processing.
// Returns true if the content is available at the endpoint, false otherwise.
func (pc PublishCheck) DoCheck() bool {
	info.Printf("Running check for UUID [%v]\n", pc.Metric.UUID)
	resp, err := http.Get(pc.Metric.endpoint.String() + pc.Metric.UUID)

	if err != nil {
		return false
	}

	check := endpointSpecificChecks[pc.Metric.config.Alias]
	if check == nil {
		warn.Printf("No check for endpoint %s.", pc.Metric.config.Alias)
		return false
	}

	return check.isCurrentOperationFinished(pc, resp)
}

func (c ContentCheck) isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool {
	// if the article was marked as deleted, operation is finished when the
	// article cannot be found anymore
	if pc.Metric.isMarkedDeleted {
		info.Printf("[%v]Marked deleted, status code [%v]", pc.Metric.UUID, response.StatusCode)
		return response.StatusCode == 404
	}

	// if not marked deleted, operation isn't finished until status is 200
	if response.StatusCode != 200 {
		return false
	}

	info.Printf("[%v]Not marked as deleted, got 200, checking PR", pc.Metric.UUID)
	// if status is 200, we check the publishReference
	// this way we can handle updates
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		warn.Printf("Cannot read response: [%s]", err.Error())
		return false
	}

	var jsonResp map[string]interface{}

	err = json.Unmarshal(data, &jsonResp)
	if err != nil {
		warn.Printf("Cannot unmarshal JSON response: [%s]", err.Error())
		return false
	}

	return jsonResp["publishReference"] == pc.Metric.tid
}

func (s S3Check) isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool {
	return response.StatusCode == 200
}

//key is the endpoint alias from the config
var endpointSpecificChecks = map[string]EndpointSpecificCheck{
	"content": ContentCheck{},
	"S3":      S3Check{},
}
