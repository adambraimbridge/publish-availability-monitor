package feeds

import (
	"net/url"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

func NewNotificationsFeed(name string, baseURL url.URL, expiry int, interval int, username string, password string, APIKey string) Feed {
	if isNotificationsPullFeed(name) {
		return newNotificationsPullFeed(name, baseURL, expiry, interval, username, password)
	} else if isNotificationsPushFeed(name) {
		return newNotificationsPushFeed(name, baseURL, expiry, interval, username, password, APIKey)
	}

	return nil
}

func isNotificationsPullFeed(feedName string) bool {
	return feedName == "notifications" ||
		feedName == "list-notifications"
}

func isNotificationsPushFeed(feedName string) bool {
	return strings.HasSuffix(feedName, "notifications-push")
}

func newNotificationsPullFeed(name string, baseURL url.URL, expiry int, interval int, username string, password string) *NotificationsPullFeed {
	feedURL := baseURL.String()

	bootstrapValues := baseURL.Query()
	bootstrapValues.Add("since", time.Now().Format(time.RFC3339))
	baseURL.RawQuery = ""

	log.Infof("constructing NotificationsPullFeed for [%s], baseUrl = [%s], bootstrapValues = [%s]", feedURL, baseURL.String(), bootstrapValues.Encode())
	return &NotificationsPullFeed{
		baseNotificationsFeed{
			name,
			nil,
			feedURL,
			username,
			password,
			expiry + 2*interval,
			make(map[string][]*Notification),
			&sync.RWMutex{},
		},
		baseURL.String(),
		bootstrapValues.Encode(),
		&sync.Mutex{},
		interval,
		nil,
		nil,
	}
}

func newNotificationsPushFeed(name string, baseURL url.URL, expiry int, interval int, username string, password string, APIKey string) *NotificationsPushFeed {
	log.Infof("constructing NotificationsPushFeed, bootstrapUrl = [%s]", baseURL.String())
	return &NotificationsPushFeed{
		baseNotificationsFeed{
			name,
			nil,
			baseURL.String(),
			username,
			password,
			expiry + 2*interval,
			make(map[string][]*Notification),
			&sync.RWMutex{},
		},
		true,
		&sync.RWMutex{},
		false,
		APIKey,
	}
}
