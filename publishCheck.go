package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// EndpointSpecificCheck is the interface which defines the URL building logic and a method
// which determines the state of the operation we are currently checking.
type EndpointSpecificCheck interface {
	buildURL(pm PublishMetric) string
	isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool
}

// ContentCheck implements the EndpointSpecificCheck interface to check operation
// status for the content endpoint.
type ContentCheck struct{}

// S3Check implements the EndpointSpecificCheck interface to check operation
// status for the S3 endpoint.
type S3Check struct{}

// NotificationsCheck implements the EndpointSpecificCheck interface to build the endpoint URL and
// to check the operation is present in the notification feed
type NotificationsCheck struct{}

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
	check := endpointSpecificChecks[pc.Metric.config.Alias]

	//TODO remove this after the notification service gets rewritten
	if pc.Metric.config.Alias == "notifications" {
		time.Sleep((120 - time.Since(pc.Metric.publishDate)) * time.Second)
	}

	resp, err := http.Get(check.buildURL(pc.Metric))
	defer resp.Body.Close()

	if err != nil {
		return false
	}

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

func (c ContentCheck) buildURL(pm PublishMetric) string {
	return pm.endpoint.String() + pm.UUID
}

func (s S3Check) isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool {
	return response.StatusCode == 200
}

func (s S3Check) buildURL(pm PublishMetric) string {
	return pm.endpoint.String() + pm.UUID
}

func (n NotificationsCheck) isCurrentOperationFinished(pc PublishCheck, response *http.Response) bool {
	if response.StatusCode != 200 {
		warn.Printf("/notifications endpoint status: [%d]", response.StatusCode)
		return false
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		warn.Printf("Cannot read response: [%s]", err.Error())
		return false
	}

	return strings.Contains(string(data), pc.Metric.UUID)
}

func (n NotificationsCheck) buildURL(pm PublishMetric) string {
	base := pm.endpoint.String()
	queryParam := url.Values{}
	//e.g. 2015-07-23T00:00:00.000Z
	since := pm.publishDate.Format(time.RFC3339Nano)
	queryParam.Add("since", since)
	return base + "?" + queryParam.Encode()
}

//key is the endpoint alias from the config
var endpointSpecificChecks = map[string]EndpointSpecificCheck{
	"content":         ContentCheck{},
	"S3":              S3Check{},
	"enrichedContent": ContentCheck{},
	"lists":           ContentCheck{},
}
