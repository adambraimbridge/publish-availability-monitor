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
	feedName          string
	httpCaller        checks.HttpCaller
	baseUrl           string
	username          string
	password          string
	expiry            int
	notifications     map[string][]*Notification
	notificationsLock *sync.RWMutex
	stopFeed          bool
	stopFeedLock      *sync.RWMutex
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

		for f.consumeFeed() {}
	}()
}

func (f *NotificationsPushFeed) Stop() {
	infoLogger.Printf("shutting down notifications push feed for %s", f.baseUrl)
	f.stopFeedLock.Lock()
	defer f.stopFeedLock.Unlock()

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
	f.notificationsLock.RLock()
	defer f.notificationsLock.RUnlock()

	var history []*Notification
	var found bool

	if history, found = f.notifications[uuid]; !found {
		history = make([]*Notification, 0)
	}

	return history
}

func (f *NotificationsPushFeed) SetHttpCaller(httpCaller checks.HttpCaller) {
	f.httpCaller = httpCaller
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

	br := bufio.NewReader(resp.Body)
	for {
		if !f.isConsuming() {
			infoLogger.Printf("stop consuming feed")
			break
		}
		f.purgeObsoleteNotifications()

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
			continue
		}

		f.storeNotifications(notifications)
	}

	return false
}

func (f *NotificationsPushFeed) parseUuidFromUrl(url string) string {
	i := strings.LastIndex(url, "/")
	return url[i+1:]
}

func (f *NotificationsPushFeed) purgeObsoleteNotifications() {
	earliest := time.Now().Add(time.Duration(-f.expiry) * time.Second).Format(time.RFC3339)
	empty := make([]string, 0)

	f.notificationsLock.Lock()
	defer f.notificationsLock.Unlock()

	for u, n := range f.notifications {
		earliestIndex := 0
		for _, e := range n {
			if strings.Compare(e.LastModified, earliest) >= 0 {
				break
			} else {
				earliestIndex++
			}
		}
		f.notifications[u] = n[earliestIndex:]

		if len(f.notifications[u]) == 0 {
			empty = append(empty, u)
		}
	}

	for _, u := range empty {
		delete(f.notifications, u)
	}
}

func (f *NotificationsPushFeed) storeNotifications(notifications []Notification) {
	f.notificationsLock.Lock()
	defer f.notificationsLock.Unlock()

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
