package feeds

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const NotificationsPull = "Notifications-Pull"

type NotificationsPullFeed struct {
	feedName          string
	httpCaller        checks.HttpCaller
	baseUrl           string
	username          string
	password          string
	sinceDate         string
	sinceDateLock     *sync.Mutex
	expiry            int
	interval          int
	ticker            *time.Ticker
	poller            chan struct{}
	notifications     map[string][]*Notification
	notificationsLock *sync.RWMutex
}

// ignore unused field (e.g. requestUrl)
type notificationsResponse struct {
	Notifications []Notification
	Links         []Link
}

func cleanupResp(resp *http.Response) {
	_, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		infoLogger.Printf("[%v]", err)
	}
	err = resp.Body.Close()
	if err != nil {
		infoLogger.Printf("[%v]", err)
	}
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

func (f *NotificationsPullFeed) FeedName() string {
	return f.feedName
}

func (f *NotificationsPullFeed) FeedType() string {
	return NotificationsPull
}

func (f *NotificationsPullFeed) SetCredentials(username string, password string) {
	f.username = username
	f.password = password
}

func (f *NotificationsPullFeed) SetHttpCaller(httpCaller checks.HttpCaller) {
	f.httpCaller = httpCaller
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
		uuid := f.parseUuidFromUrl(n.ID)
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

func (f *NotificationsPullFeed) purgeObsoleteNotifications() {
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

func (f *NotificationsPullFeed) buildNotificationsURL() string {
	q := url.Values{}
	q.Add("since", f.sinceDate)

	return f.baseUrl + "?" + q.Encode()
}

func (f *NotificationsPullFeed) parseUuidFromUrl(url string) string {
	i := strings.LastIndex(url, "/")
	return url[i+1:]
}

func (f *NotificationsPullFeed) NotificationsFor(uuid string) []*Notification {
	var history []*Notification
	var found bool

	f.notificationsLock.RLock()
	defer f.notificationsLock.RUnlock()

	if history, found = f.notifications[uuid]; !found {
		history = make([]*Notification, 0)
	}

	return history
}
