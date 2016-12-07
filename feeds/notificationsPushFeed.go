package feeds

import (
	"bufio"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const NotificationsPush = "Notifications-Push"

type NotificationsPushFeed struct {
	baseNotificationsFeed
	stopFeed     bool
	stopFeedLock *sync.RWMutex
	connected    bool
}

func (f *NotificationsPushFeed) Start() {
	infoLogger.Printf("starting notifications-push feed from %v", f.baseUrl)
	f.stopFeedLock.Lock()
	defer f.stopFeedLock.Unlock()

	f.stopFeed = false
	go func() {
		if f.httpCaller == nil {
			f.httpCaller = checks.NewHttpCaller(0)
		}

		for f.consumeFeed() {
			time.Sleep(500 * time.Millisecond)
			infoLogger.Println("Disconnected from Push feed! Attempting to reconnect.")
		}
	}()
}

func (f *NotificationsPushFeed) Stop() {
	infoLogger.Printf("shutting down notifications push feed for %s", f.baseUrl)
	f.stopFeedLock.Lock()
	defer f.stopFeedLock.Unlock()

	f.stopFeed = true
}

func (f *NotificationsPushFeed) FeedType() string {
	return NotificationsPush
}

func (f *NotificationsPushFeed) IsConnected() bool {
	return f.connected
}

func (f *NotificationsPushFeed) isConsuming() bool {
	f.stopFeedLock.RLock()
	defer f.stopFeedLock.RUnlock()

	return !f.stopFeed
}

func (f *NotificationsPushFeed) consumeFeed() bool {
	resp, err := f.httpCaller.DoCall(f.baseUrl, f.username, f.password)

	if err != nil {
		infoLogger.Printf("Sending request: [%v]", err)
		return f.isConsuming()
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		infoLogger.Printf("Received invalid statusCode: [%v]", resp.StatusCode)
		return f.isConsuming()
	}

	infoLogger.Println("Reconnected to push feed!")
	f.connected = true
	defer func() { f.connected = false }()

	br := bufio.NewReader(resp.Body)
	for {
		if !f.isConsuming() {
			infoLogger.Printf("stop consuming feed")
			break
		}
		f.purgeObsoleteNotifications()

		event, err := br.ReadString('\n')
		if err != nil {
			infoLogger.Printf("Error: [%v] - disconnected from push feed!", err)
			return f.isConsuming()
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
			continue
		}

		f.storeNotifications(notifications)
	}

	return false
}

func (f *NotificationsPushFeed) storeNotifications(notifications []Notification) {
	f.notificationsLock.Lock()
	defer f.notificationsLock.Unlock()

	for _, n := range notifications {
		uuid := parseUuidFromUrl(n.ID)
		var history []*Notification
		var found bool
		if history, found = f.notifications[uuid]; !found {
			history = make([]*Notification, 0)
		}

		history = append(history, &n)
		f.notifications[uuid] = history
	}
}
