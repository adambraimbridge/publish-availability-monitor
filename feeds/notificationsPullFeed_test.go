package feeds

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	t             *testing.T
	authUser      string
	authPass      string
	txIdPrefix    string
	mockResponses []*http.Response
	current       int
}

// returns the mock responses of testHTTPCaller in order
func (t *testHTTPCaller) DoCall(url string, username string, password string, txId string) (*http.Response, error) {
	if t.authUser != username || t.authPass != password {
		return buildResponse(401, `{message: "Not authenticated"}`), nil
	}

	if t.txIdPrefix != "" {
		assert.True(t.t, strings.HasPrefix(txId, t.txIdPrefix), "transaction id should start with "+t.txIdPrefix)
		timestamp := txId[len(t.txIdPrefix):]
		_, err := time.Parse(time.RFC3339, timestamp)
		assert.Nil(t.t, err, "transaction id suffix did not parse as a timestamp")
	}

	response := t.mockResponses[t.current]
	t.current = (t.current + 1) % len(t.mockResponses)
	return response, nil
}

func (t *testHTTPCaller) DoCallWithEntity(httpMethod string, url string, username string, password string, txId string, contentType string, entity io.Reader) (*http.Response, error) {
	return nil, nil
}

// builds testHTTPCaller with the given mocked responses in the provided order
func mockHTTPCaller(t *testing.T, txIdPrefix string, responses ...*http.Response) checks.HttpCaller {
	return &testHTTPCaller{t: t, txIdPrefix: txIdPrefix, mockResponses: responses}
}

// builds testHTTPCaller with the given mocked responses in the provided order
func mockAuthenticatedHTTPCaller(t *testing.T, txIdPrefix string, username string, password string, responses ...*http.Response) checks.HttpCaller {
	return &testHTTPCaller{t: t, txIdPrefix: txIdPrefix, authUser: username, authPass: password, mockResponses: responses}
}

// this is necessary to be able to build an http.Response
// the body has to be a ReadCloser
type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func mockNotificationFor(uuid string, publishRef string, lastModified time.Time) string {
	return fmt.Sprintf(`{
						"type": "http://www.ft.com/thing/ThingChangeType/UPDATE",
						"id": "http://www.ft.com/thing/%v",
						"apiUrl": "http://api.ft.com/content/%v",
						"publishReference": "%v",
						"lastModified": "%v"
					}`, uuid, uuid, publishRef, lastModified.Format(time.RFC3339))
}

func mockNotificationsResponseFor(sinceDate string, notifications string, nextSinceDate string) string {
	return fmt.Sprintf(`{
			"requestUrl": "http://api.ft.com/content/notifications?since=%v",
			"notifications": [
					%v
			],
			"links": [
					{
						"href": "http://api.ft.com/content/notifications?since=%v",
						"rel": "next"
					}
			]
		}`, sinceDate, notifications, nextSinceDate)
}

func TestBuiltNotificationsPullURLIsCorrect(t *testing.T) {
	baseUrl, _ := url.Parse("http://www.example.org?type=all")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 10, 1, "", "")

	actual := f.(*NotificationsPullFeed).buildNotificationsURL()

	uri, err := url.Parse(actual)
	assert.NoError(t, err)
	assert.Equal(t, "all", uri.Query().Get("type"))
	assert.Equal(t, sinceDate, uri.Query().Get("since"))
	assert.Equal(t, "www.example.org", uri.Host)
	assert.Equal(t, "http", uri.Scheme)
}

func TestNotificationsArePolled(t *testing.T) {
	uuid := "1cb14245-5185-4ed5-9188-4d2a86085599"
	publishRef := "tid_0123wxyz"
	lastModified := time.Now()
	notifications := mockNotificationsResponseFor("2016-10-28T15:00:00.000Z",
		mockNotificationFor(uuid, publishRef, lastModified),
		"2016-10-28T16:00:00.000Z")

	httpCaller := mockHTTPCaller(t, "tid_pam_notifications_pull_", buildResponse(200, notifications))

	baseUrl, _ := url.Parse("http://www.example.org?type=all")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 10, 1, "", "")

	f.(*NotificationsPullFeed).SetHttpCaller(httpCaller)
	f.Start()
	defer f.Stop()

	time.Sleep(time.Duration(1200) * time.Millisecond)

	response := f.NotificationsFor(uuid)
	assert.Len(t, response, 1, "notifications for item")
	assert.Equal(t, publishRef, response[0].PublishReference, "publish ref")
}

func TestNotificationsForReturnsEmptyIfNotFound(t *testing.T) {
	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 10, 1, "", "")

	response := f.NotificationsFor("1cb14245-5185-4ed5-9188-4d2a86085599")
	assert.Len(t, response, 0, "notifications for item")
}

func TestNotificationsForReturnsAllMatches(t *testing.T) {
	uuid := "1cb14245-5185-4ed5-9188-4d2a86085599"
	publishRef1 := "tid_0123wxyz"
	lastModified1 := time.Now().Add(time.Duration(-1) * time.Second)
	notifications1 := mockNotificationsResponseFor("2016-10-28T15:00:00.000Z",
		mockNotificationFor(uuid, publishRef1, lastModified1),
		"2016-10-28T15:10:00.000Z")

	publishRef2 := "tid_0123abcd"
	lastModified2 := time.Now()
	notifications2 := mockNotificationsResponseFor("2016-10-28T15:10:00.000Z",
		mockNotificationFor(uuid, publishRef2, lastModified2),
		"2016-10-28T15:20:00.000Z")

	httpCaller := mockHTTPCaller(t, "tid_pam_notifications_pull_", buildResponse(200, notifications1), buildResponse(200, notifications2))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 10, 1, "", "")
	f.(*NotificationsPullFeed).SetHttpCaller(httpCaller)
	f.Start()
	defer f.Stop()
	time.Sleep(time.Duration(2200) * time.Millisecond)

	response := f.NotificationsFor(uuid)
	assert.Len(t, response, 2, "notifications for item")
	assert.Equal(t, publishRef1, response[0].PublishReference, "first publish ref")
	assert.Equal(t, publishRef2, response[1].PublishReference, "second publish ref")
}

func TestNotificationsPollingContinuesAfterErrorResponse(t *testing.T) {
	uuid := "1cb14245-5185-4ed5-9188-4d2a86085599"
	publishRef := "tid_0123wxyz"
	lastModified := time.Now()
	notifications := mockNotificationsResponseFor("2016-10-28T15:00:00.000Z",
		mockNotificationFor(uuid, publishRef, lastModified),
		"2016-10-28T16:00:00.000Z")

	httpCaller := mockHTTPCaller(t, "tid_pam_notifications_pull_", buildResponse(500, ""), buildResponse(200, notifications))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 10, 1, "", "")
	f.(*NotificationsPullFeed).SetHttpCaller(httpCaller)
	f.Start()
	defer f.Stop()
	time.Sleep(time.Duration(2200) * time.Millisecond)

	response := f.NotificationsFor(uuid)
	assert.Len(t, response, 1, "notifications for item")
	assert.Equal(t, publishRef, response[0].PublishReference, "publish ref")
}

func TestNotificationsArePurged(t *testing.T) {
	uuid := "1cb14245-5185-4ed5-9188-4d2a86085599"
	publishRef := "tid_0123wxyz"
	lastModified := time.Now().Add(time.Duration(-2) * time.Second)
	notifications := mockNotificationsResponseFor("2016-10-28T15:00:00.000Z",
		mockNotificationFor(uuid, publishRef, lastModified),
		"2016-10-28T16:00:00.000Z")

	httpCaller := mockHTTPCaller(t, "tid_pam_notifications_pull_", buildResponse(200, notifications))

	baseUrl, _ := url.Parse("http://www.example.org")
	sinceDate := "2016-10-28T15:00:00.000Z"
	f := NewNotificationsFeed("notifications", baseUrl, sinceDate, 1, 1, "", "")
	f.(*NotificationsPullFeed).SetHttpCaller(httpCaller)
	f.Start()
	defer f.Stop()

	time.Sleep(time.Duration(1200) * time.Millisecond)

	response := f.NotificationsFor(uuid)
	assert.Len(t, response, 1, "notifications for item")
	assert.Equal(t, publishRef, response[0].PublishReference, "publish ref")

	time.Sleep(time.Duration(1) * time.Second)
	response = f.NotificationsFor(uuid)
	assert.Len(t, response, 0, "notifications for item")
}
