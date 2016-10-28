package feeds

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/stretchr/testify/assert"
)

func buildResponse(statusCode int, content string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       nopCloser{bytes.NewBuffer([]byte(content))},
	}
}

// mock httpCaller implementation
type testHTTPCaller struct {
	authUser      string
	authPass      string
	mockResponses []*http.Response
	current       int
}

// returns the mock responses of testHTTPCaller in order
func (t *testHTTPCaller) DoCall(url string, username string, password string) (*http.Response, error) {
	if t.authUser != username || t.authPass != password {
		return buildResponse(401, `{message: "Not authenticated"}`), nil
	}

	response := t.mockResponses[t.current]
	t.current = (t.current + 1) % len(t.mockResponses)
	return response, nil
}

// builds testHTTPCaller with the given mocked responses in the provided order
func mockHTTPCaller(responses ...*http.Response) checks.HttpCaller {
	return &testHTTPCaller{mockResponses: responses}
}

// builds testHTTPCaller with the given mocked responses in the provided order
func mockAuthenticatedHTTPCaller(username string, password string, responses ...*http.Response) checks.HttpCaller {
	return &testHTTPCaller{authUser: username, authPass: password, mockResponses: responses}
}

// this is necessary to be able to build an http.Response
// the body has to be a ReadCloser
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestNotificationsArePolled(t *testing.T) {
	notifications := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2016-10-28T15:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz",
						"lastModified": "2016-10-28T15:16:47.391Z"
					}
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2016-10-28T16:00:00.000Z",
						"rel": "next"
					}
			]
		}`

	httpCaller := mockHTTPCaller(buildResponse(200, notifications))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsPullFeed(httpCaller, baseUrl, sinceDate, 1, "", "")
	f.Start()
	time.Sleep(time.Duration(1200) * time.Millisecond)
	f.Stop()

	response := f.NotificationsFor("1cb14245-5185-4ed5-9188-4d2a86085599")
	assert.Len(t, response, 1, "notifications for item")
	assert.Equal(t, "tid_0123wxyz", response[0].PublishReference, "publish ref")
}

func TestNotificationsForReturnsEmptyIfNotFound(t *testing.T) {
	httpCaller := mockHTTPCaller(buildResponse(200, ""))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsPullFeed(httpCaller, baseUrl, sinceDate, 1, "", "")

	response := f.NotificationsFor("1cb14245-5185-4ed5-9188-4d2a86085599")
	assert.Len(t, response, 0, "notifications for item")
}

func TestNotificationsForReturnsAllMatches(t *testing.T) {
	notifications1 := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2016-10-28T15:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz",
						"lastModified": "2016-10-28T15:06:47.391Z"
					}
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2016-10-28T15:10:00.000Z",
						"rel": "next"
					}
			]
		}`

	notifications2 := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2016-10-28T15:10:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123abcd",
						"lastModified": "2016-10-28T15:16:47.391Z"
					}
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2016-10-28T15:20:00.000Z",
						"rel": "next"
					}
			]
		}`
	httpCaller := mockHTTPCaller(buildResponse(200, notifications1), buildResponse(200, notifications2))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsPullFeed(httpCaller, baseUrl, sinceDate, 1, "", "")
	f.Start()
	time.Sleep(time.Duration(2200) * time.Millisecond)
	f.Stop()

	response := f.NotificationsFor("1cb14245-5185-4ed5-9188-4d2a86085599")
	assert.Len(t, response, 2, "notifications for item")
	assert.Equal(t, "tid_0123wxyz", response[0].PublishReference, "first publish ref")
	assert.Equal(t, "tid_0123abcd", response[1].PublishReference, "second publish ref")
}

func TestNotificationsPollingContinuesAfterErrorResponse(t *testing.T) {
	notifications := `{
			"requestUrl": "http://api.ft.com/content/notifications?since=2016-10-28T15:00:00.000Z",
			"notifications": [
					{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/1cb14245-5185-4ed5-9188-4d2a86085599",
						"apiUrl": "http://api.ft.com/content/1cb14245-5185-4ed5-9188-4d2a86085599",
						"publishReference": "tid_0123wxyz",
						"lastModified": "2016-10-28T15:16:47.391Z"
					}
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=2016-10-28T16:00:00.000Z",
						"rel": "next"
					}
			]
		}`

	httpCaller := mockHTTPCaller(buildResponse(500, ""), buildResponse(200, notifications))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsPullFeed(httpCaller, baseUrl, sinceDate, 1, "", "")
	f.Start()
	time.Sleep(time.Duration(2200) * time.Millisecond)
	f.Stop()

	response := f.NotificationsFor("1cb14245-5185-4ed5-9188-4d2a86085599")
	assert.Len(t, response, 1, "notifications for item")
	assert.Equal(t, "tid_0123wxyz", response[0].PublishReference, "publish ref")
}
