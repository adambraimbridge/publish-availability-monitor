package feeds

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const NotificationsPull = "Notifications-Pull"

type NotificationsPullFeed struct {
	httpCaller    checks.HttpCaller
	baseUrl       string
	username      string
	password      string
	sinceDate     string
	expiry        int
	interval      int
	ticker        *time.Ticker
	poller        chan struct{}
	notifications map[string][]*Notification
}

// ignore unused field (e.g. requestUrl)
type notificationsResponse struct {
	Notifications []Notification
	Links         []Link
}

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  - ", logPattern)
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

func NewNotificationsPullFeed(httpCaller checks.HttpCaller, baseUrl *url.URL, sinceDate string, expiry int, interval int, username string, password string) *NotificationsPullFeed {
	return &NotificationsPullFeed{httpCaller: httpCaller,
		baseUrl:       baseUrl.String(),
		sinceDate:     sinceDate,
		expiry:        expiry + 2*interval,
		interval:      interval,
		username:      username,
		password:      password,
		notifications: make(map[string][]*Notification)}
}

func (f *NotificationsPullFeed) Start() {
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

func (f *NotificationsPullFeed) Name() string {
	return NotificationsPull
}

func (f *NotificationsPullFeed) SetCredentials(username string, password string) {
	f.username = username
	f.password = password
}

func (f *NotificationsPullFeed) pollNotificationsFeed() {
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

	if history, found = f.notifications[uuid]; !found {
		history = make([]*Notification, 0)
	}

	return history
}
