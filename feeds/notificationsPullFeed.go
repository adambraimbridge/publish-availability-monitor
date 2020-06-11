package feeds

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	log "github.com/Sirupsen/logrus"
)

const NotificationsPull = "Notifications-Pull"

type NotificationsPullFeed struct {
	baseNotificationsFeed
	notificationsURL         string
	notificationsQueryString string
	notificationsURLLock     *sync.Mutex
	interval                 int
	ticker                   *time.Ticker
	poller                   chan struct{}
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
	log.Infof("shutting down notifications pull feed for %s", f.baseURL)
	close(f.poller)
}

func (f *NotificationsPullFeed) FeedType() string {
	return NotificationsPull
}

func (f *NotificationsPullFeed) pollNotificationsFeed() {
	f.notificationsURLLock.Lock()
	defer f.notificationsURLLock.Unlock()

	txID := f.buildNotificationsTxID()
	notificationsURL := f.notificationsURL + "?" + f.notificationsQueryString
	resp, err := f.httpCaller.DoCall(checks.Config{URL: notificationsURL, Username: f.username, Password: f.password, TxID: txID}) //nolint:bodyclose

	if err != nil {
		log.WithField("transaction_id", txID).WithError(err).Errorf("error calling notifications %s", notificationsURL)
		return
	}
	defer cleanupResp(resp)

	if resp.StatusCode != 200 {
		log.WithField("transaction_id", txID).Errorf("Notifications [%s] status NOT OK: [%d]", notificationsURL, resp.StatusCode)
		return
	}

	var notifications notificationsResponse
	err = json.NewDecoder(resp.Body).Decode(&notifications)
	if err != nil {
		log.WithField("transaction_id", txID).Errorf("Cannot decode json response: [%s]", err.Error())
		return
	}

	f.notificationsLock.Lock()
	defer f.notificationsLock.Unlock()

	for _, v := range notifications.Notifications {
		n := v
		uuid := parseUuidFromUrl(n.ID)
		var history []*Notification
		var found bool
		if history, found = f.notifications[uuid]; !found {
			history = make([]*Notification, 0)
		}

		history = append(history, &n)
		f.notifications[uuid] = history
	}

	nextPageURL, err := url.Parse(notifications.Links[0].Href)
	if err != nil {
		log.Errorf("unparseable next url: [%s]", notifications.Links[0].Href)
		return // and hope that a retry will fix this
	}

	f.notificationsQueryString = nextPageURL.RawQuery
}

func (f *NotificationsPullFeed) buildNotificationsTxID() string {
	return "tid_pam_notifications_pull_" + time.Now().Format(time.RFC3339)
}
