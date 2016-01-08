package main

import (
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestCheckBatchOfNotifications_ResponseBatchOfNotificationsIsEmpty_NotFinished(t *testing.T) {
	testResponse := `{
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
		mockHTTPCaller(buildResponse(200, testResponse)),
	}
	if finished, _ := notificationsCheck.checkBatchOfNotifications("dummy-url", "tid_0123wxyZ"); finished {
		t.Error("Expected failure")
	}
}

func TestCheckBatchOfNotifications_ResponseDoesNotContainTID_NotFinished(t *testing.T) {
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

func TestCheckBatchOfNotifications_ResponseDoesContainTID_Finished(t *testing.T) {
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

func TestIsCurrentOperationFinished_FirstBatchOfNotificationsContainsTID_Finished(t *testing.T) {
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

	if finished, _ := notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(testTID).build()); !finished {
		t.Error("Expected success")
	}
}

func TestIsCurrentOperationFinished_FirstBatchOfNotificationsDoesNotContainTIDSecondBatchIsEmpty_NotFinished(t *testing.T) {
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
	if finished, _ := notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID("tid_0123wxyZ").build()); finished {
		t.Error("Expected failure")
	}
}

func TestIsCurrentOperationFinished_FirstBatchOfNotificationsDoesNotContainTIDButSecondDoes_Finished(t *testing.T) {
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

	if finished, _ := notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTID).build()); !finished {
		t.Error("Expected success")
	}
}

func TestNotificationsBuildURL_Success(test *testing.T) {
	publishDate, err := time.Parse(time.RFC3339Nano, "2015-10-21T14:22:06.271Z")
	if err != nil {
		test.Errorf("Test data error: [%v]", err)
	}
	pm := newPublishMetricBuilder().withEndpoint("http://notifications-endpoint:8080/content/notifications").withPublishDate(publishDate).build()

	expected := "http://notifications-endpoint:8080/content/notifications?since=2015-10-21T14:22:06.271Z"

	actual, err := url.QueryUnescape(buildNotificationsURL(pm))
	if err != nil {
		test.Errorf("Expected success. Found error: [%v]", err)
	}
	if actual != expected {
		test.Errorf("Expected success.\nActual: [%s]\nExpected: [%s]", actual, expected)
	}
}

func TestAdjustNextNotificationsURL_CurrentHostAndPortDiffers_Success(t *testing.T) {
	current := "http://ftapp14927-lvpr-uk-int:8080/content/notifications?since=2015-12-15T00:00:00.000Z"
	next := "http://int.api.ft.com/content/notifications?since=2015-12-15T11:53:17.508Z"

	expected := "http://ftapp14927-lvpr-uk-int:8080/content/notifications?since=2015-12-15T11:53:17.508Z"

	actual, err := adjustNextNotificationsURL(current, next)
	if err != nil {
		t.Errorf("Expected success. Found error: [%v]", err)
	}

	if actual != expected {
		t.Error("Expected success")
	}
}
