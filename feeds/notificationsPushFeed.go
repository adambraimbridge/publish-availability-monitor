package feeds

import (
	"bufio"
	"encoding/json"
	"strings"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const NotificationsPush = "Notifications-Push"

type NotificationsPushFeed struct {
	feedName      string
	httpCaller    checks.HttpCaller
	baseUrl       string
	username      string
	password      string
	notifications map[string][]*Notification
	stopFeed      bool
}

func (f *NotificationsPushFeed) Start() {
	infoLogger.Printf("starting notifications-push feed from %v", f.baseUrl)
	f.stopFeed = false
	go f.consumeFeed()
}

func (f *NotificationsPushFeed) Stop() {
	infoLogger.Printf("shutting down notifications pull feed for %s", f.baseUrl)
	f.stopFeed = true
}

func (f *NotificationsPushFeed) FeedName() string {
	return f.feedName
}

func (f *NotificationsPushFeed) FeedType() string {
	return NotificationsPush
}

func (f *NotificationsPushFeed) SetCredentials(username string, password string) {
	f.username = username
	f.password = password
}

func (f *NotificationsPushFeed) NotificationsFor(uuid string) []*Notification {
	infoLogger.Printf("checking notifications for %v", uuid)
	var history []*Notification
	var found bool

	if history, found = f.notifications[uuid]; !found {
		history = make([]*Notification, 0)
	}

	infoLogger.Printf("notifications for %v: %v", uuid, history)
	return history
}

func (f *NotificationsPushFeed) consumeFeed() {
	f.httpCaller = checks.NewHttpCaller(0)
	resp, err := f.httpCaller.DoCall(f.baseUrl, f.username, f.password)

	if err != nil {
		infoLogger.Fatalf("Sending request: [%v]", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		infoLogger.Fatalf("Received invalid statusCode: [%v]", resp.StatusCode)
	}

	br := bufio.NewReader(resp.Body)
	for {
		if f.stopFeed {
			infoLogger.Printf("stop consuming feed")
			break
		}

		event, err := br.ReadString('\n')
		if err != nil {
			infoLogger.Printf("Error: [%v]", err)
			panic("foo")
			continue
		}

		trimmed := strings.TrimSpace(event)
		if trimmed == "" {
			continue
		}

		data := strings.TrimPrefix(trimmed, "data: ")
		var notifications []Notification
		err = json.Unmarshal([]byte(data), &notifications)
		if err != nil {
			infoLogger.Printf("Error: [%v]. \n", err)
			continue
		}

		if len(notifications) == 0 {
			infoLogger.Print("Received 'heartbeat' event")
			continue
		}

		infoLogger.Printf("Received notifications: [%v]", notifications)

		for _, n := range notifications {
			uuid := f.parseUuidFromUrl(n.ID)
			var history []*Notification
			var found bool
			if history, found = f.notifications[uuid]; !found {
				history = make([]*Notification, 0)
			}

			history = append(history, &n)
			f.notifications[uuid] = history
		}
	}
}

func (f *NotificationsPushFeed) parseUuidFromUrl(url string) string {
	i := strings.LastIndex(url, "/")
	return url[i+1:]
}
