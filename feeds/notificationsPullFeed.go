package feeds

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const NotificationsPull = "Notifications-Pull"

type NotificationsPullFeed struct {
	baseNotificationsFeed
	sinceDate     string
	sinceDateLock *sync.Mutex
	interval      int
	ticker        *time.Ticker
	poller        chan struct{}
}

// ignore unused field (e.g. requestUrl)
type notificationsResponse struct {
	Notifications []Notification
	Links         []Link
}

func (f *NotificationsPullFeed) Start() {
	if f.httpCaller == nil {
		f.httpCaller = checks.NewHttpCaller(10)
	}

	f.ticker = time.NewTicker(time.Duration(f.interval) * time.Second)
	f.poller = make(chan struct{})
	go func() {
		for {
			select {
			case <-f.ticker.C:
				go func() {
					f.pollNotificationsFeed()
					f.purgeObsoleteNotifications()
				}()
			case <-f.poller:
				f.ticker.Stop()
				return
			}
		}
	}()
}

func (f *NotificationsPullFeed) Stop() {
	infoLogger.Printf("shutting down notifications pull feed for %s", f.baseUrl)
	close(f.poller)
}

func (f *NotificationsPullFeed) FeedType() string {
	return NotificationsPull
}

func (f *NotificationsPullFeed) pollNotificationsFeed() {
	f.sinceDateLock.Lock()
	defer f.sinceDateLock.Unlock()

	notificationsUrl := f.buildNotificationsURL()
	resp, err := f.httpCaller.DoCall(notificationsUrl, f.username, f.password)

	if err != nil {
		infoLogger.Printf("error calling notifications %s", notificationsUrl)
		return
	}
	defer cleanupResp(resp)

	if resp.StatusCode != 200 {
		infoLogger.Printf("Notifications [%s] status NOT OK: [%d]", notificationsUrl, resp.StatusCode)
		return
	}

	var notifications notificationsResponse
	err = json.NewDecoder(resp.Body).Decode(&notifications)
	if err != nil {
		infoLogger.Printf("Cannot decode json response: [%s]", err.Error())
		return
	}

	f.notificationsLock.Lock()
	defer f.notificationsLock.Unlock()

	for _, n := range notifications.Notifications {
		uuid := parseUuidFromUrl(n.ID)
		var history []*Notification
		var found bool
		if history, found = f.notifications[uuid]; !found {
			history = make([]*Notification, 0)
		}

		history = append(history, &n)
		f.notifications[uuid] = history
	}

	nextPageUrl, _ := url.Parse(notifications.Links[0].Href)
	f.sinceDate = nextPageUrl.Query().Get("since")
}

func (f *NotificationsPullFeed) buildNotificationsURL() string {
	q := url.Values{}
	q.Add("since", f.sinceDate)

	return f.baseUrl + "?" + q.Encode()
}
