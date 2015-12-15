package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestIsCurrentOperationFinished_S3Check_Finished(t *testing.T) {
	s3Check := &S3Check{
		mockHTTPCaller(buildResponse(200, "imagebytes")),
	}
	if !s3Check.isCurrentOperationFinished(PublishMetric{}) {
		t.Errorf("Expected: true. Actual: false")
	}
}

func TestIsCurrentOperationFinished_S3Check_Empty(t *testing.T) {
	s3Check := &S3Check{
		mockHTTPCaller(buildResponse(200, "")),
	}
	if s3Check.isCurrentOperationFinished(PublishMetric{}) {
		t.Errorf("Expected: false. Actual: true")
	}
}

func TestIsCurrentOperationFinished_S3Check_NotFinished(t *testing.T) {
	s3Check := &S3Check{
		mockHTTPCaller(buildResponse(404, "")),
	}
	if s3Check.isCurrentOperationFinished(PublishMetric{}) {
		t.Errorf("Expected: false. Actual: True")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_InvalidContent(t *testing.T) {
	testResponse := `{ "uuid" : "1234-1234"`
	contentCheck := &ContentCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}

	if contentCheck.isCurrentOperationFinished(PublishMetric{}) {
		t.Errorf("Expected error.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_Finished(t *testing.T) {
	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)
	contentCheck := &ContentCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}

	if !contentCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTid).build()) {
		t.Error("Expected success.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_NotFinished(t *testing.T) {
	currentTid := "tid_1234"
	testResponse := `{ "uuid" : "1234-1234", "publishReference" : "tid_1235"}`
	contentCheck := &ContentCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}

	if contentCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTid).build()) {
		t.Error("Expected failure.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_MarkedDeleted_Finished(t *testing.T) {
	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)
	contentCheck := &ContentCheck{
		mockHTTPCaller(buildResponse(404, testResponse)),
	}

	if !contentCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTid).withMarkedDeleted(true).build()) {
		t.Error("Expected success.")
	}
}

func TestIsCurrentOperationFinished_ContentCheck_MarkedDeleted_NotFinished(t *testing.T) {
	currentTid := "tid_1234"
	testResponse := fmt.Sprintf(`{ "uuid" : "1234-1234", "publishReference" : "%s"}`, currentTid)
	contentCheck := &ContentCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}

	if contentCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTid).withMarkedDeleted(true).build()) {
		t.Error("Expected failure.")
	}
}

func TestDoSingleNotificationsPageCheck_ResponseDoesNotContainTID_NotFinished(t *testing.T) {
	testResponse := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T00:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz"
					}
				],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`
	notificationsCheck := &NotificationsCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}
	if finished, _ := notificationsCheck.checkBatchOfNotifications("dummy-url", "tid_0123wxyZ"); finished {
		t.Error("Expected failure")
	}
}

func TestDoSingleNotificationsPageCheck_ResponseDoesContainTID_Finished(t *testing.T) {
	currentTID := "tid_0123wxyZ"
	testResponse := fmt.Sprintf(`{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T00:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "%s"
					}
				],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`, currentTID)
	notificationsCheck := &NotificationsCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}
	if finished, _ := notificationsCheck.checkBatchOfNotifications("dummy-url", currentTID); !finished {
		t.Error("Expected success")
	}
}

func TestIsCurrentOperationFinished_NotificationsCheck_FirstNotificationsPageContainsTID_Finished(t *testing.T) {
	testTID := "tid_0123wxyz"
	testResponse := fmt.Sprintf(
		`{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T00:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "%s"
					}
				],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`, testTID)

	notificationsCheck := &NotificationsCheck{
		mockHTTPCaller(buildResponse(200, testResponse)),
	}

	if !notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(testTID).build()) {
		t.Error("Expected success")
	}
}

func TestIsCurrentOperationFinished_NotificationsCheck_FirstNotificationsPageDoesNotContainTIDSecondIsEmptyList_NotFinished(t *testing.T) {
	testResponse1 := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T00:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz"
					}
				],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`
	testResponse2 := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
			"notifications": [],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`
	notificationsCheck := &NotificationsCheck{
		mockHTTPCaller(buildResponse(200, testResponse1), buildResponse(200, testResponse2)),
	}
	if notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID("tid_0123wxyZ").build()) {
		t.Error("Expected failure")
	}
}

func TestIsCurrentOperationFinished_NotificationsCheck_FirstPageDoesNotContainTIDButSecondDoes_Finished(t *testing.T) {
	currentTID := "tid_0123wxyZ"
	testResponse1 := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T00:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz"
					}
				],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
						"rel": "next"
					}
			]
		}`
	testResponse2 := fmt.Sprintf(`{
			"requestUrl": "http://api.ft.com/content/notifications?since=2015-11-09T14:09:08.705Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "%s"
					}
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2015-11-09T14:10:08.500Z",
						"rel": "next"
					}
			]
		}`, currentTID)

	notificationsCheck := &NotificationsCheck{
		mockHTTPCaller(buildResponse(200, testResponse1), buildResponse(200, testResponse2)),
	}

	if !notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTID).build()) {
		t.Error("Expected success")
	}
}

func TestNotificationsBuildURL_SinceQueryParamsCorrectlyParsed(test *testing.T) {
	publishDate, err := time.Parse(time.RFC3339Nano, "2015-10-21T14:22:06.270Z")
	if err != nil {
		test.Errorf("Error in test data: [%v]", err)
	}

	pm := newPublishMetricBuilder().withEndpoint("http://notifications-endpoint:8080/content/notifications").withPublishDate(publishDate).build()

	builtURL, err := url.Parse(buildNotificationsURL(pm))
	if err != nil {
		test.Errorf("Cannot parse built URL: [%s].", err.Error())
	}

	queryParams := builtURL.Query()
	since := queryParams.Get("since")
	if since == "" {
		test.Errorf("Missing 'since' query parameter.")
	}

	t, err := time.Parse(time.RFC3339Nano, since)
	if err != nil {
		test.Errorf("Cannot parse param value: [%s]. Error: [%s]", since, err.Error())
	}

	if !t.Equal(publishDate) {
		test.Errorf("Expected timestamp: [%v]. Actual: [%v].", publishDate, t)
	}
}

type publishMetricBuilder interface {
	withUUID(string) publishMetricBuilder
	withEndpoint(string) publishMetricBuilder
	withTID(string) publishMetricBuilder
	withMarkedDeleted(bool) publishMetricBuilder
	withPublishDate(time.Time) publishMetricBuilder
	build() PublishMetric
}

//PublishMetricBuilder implementation
type pmBuilder struct {
	UUID          string
	endpoint      url.URL
	tid           string
	markedDeleted bool
	publishDate   time.Time
}

func (b *pmBuilder) withUUID(uuid string) publishMetricBuilder {
	b.UUID = uuid
	return b
}

func (b *pmBuilder) withEndpoint(endpoint string) publishMetricBuilder {
	e, _ := url.Parse(endpoint)
	b.endpoint = *e
	return b
}

func (b *pmBuilder) withTID(tid string) publishMetricBuilder {
	b.tid = tid
	return b
}

func (b *pmBuilder) withMarkedDeleted(markedDeleted bool) publishMetricBuilder {
	b.markedDeleted = markedDeleted
	return b
}

func (b *pmBuilder) withPublishDate(publishDate time.Time) publishMetricBuilder {
	b.publishDate = publishDate
	return b
}

func (b *pmBuilder) build() PublishMetric {
	return PublishMetric{
		UUID:            b.UUID,
		endpoint:        b.endpoint,
		tid:             b.tid,
		isMarkedDeleted: b.markedDeleted,
		publishDate:     b.publishDate,
	}
}

func newPublishMetricBuilder() publishMetricBuilder {
	return &pmBuilder{}
}

func buildResponse(statusCode int, content string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       nopCloser{bytes.NewBuffer([]byte(content))},
	}
}

//mock HttpCaller implementation
type testHTTPCaller struct {
	mockResponses []*http.Response
	current       int
}

func (t *testHTTPCaller) doCall(url string) (*http.Response, error) {
	response := t.mockResponses[t.current]
	t.current = (t.current + 1) % len(t.mockResponses)
	return response, nil
}

//builds testHTTPCaller with the given mocked responses in the provided order
func mockHTTPCaller(responses ...*http.Response) httpCaller {
	return &testHTTPCaller{responses, 0}
}

// this is necessary to be able to build an http.Response
// the body has to be a ReadCloser
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }
