package main

import (
	"fmt"
	"testing"
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

	if !notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(testTID).build()) {
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
	if notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID("tid_0123wxyZ").build()) {
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

	if !notificationsCheck.isCurrentOperationFinished(newPublishMetricBuilder().withTID(currentTID).build()) {
		t.Error("Expected success")
	}
}
