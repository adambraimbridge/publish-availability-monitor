package feeds

import (
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  - ", logPattern)
}

func NewNotificationsFeed(name string, baseUrl *url.URL, sinceDate string, expiry int, interval int, username string, password string) Feed {
	if isNotificationsPullFeed(name) {
		return &NotificationsPullFeed{
			name,
			nil,
			baseUrl.String(),
			username,
			password,
			sinceDate,
			&sync.Mutex{},
			expiry + 2*interval,
			interval,
			nil,
			nil,
			make(map[string][]*Notification),
			&sync.RWMutex{},
		}
	} else if isNotificationsPushFeed(name) {
		return &NotificationsPushFeed{
			name,
			nil,
			baseUrl.String(),
			username,
			password,
			expiry + 2*interval,
			make(map[string][]*Notification),
			&sync.RWMutex{},
			true,
			&sync.RWMutex{},
		}
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
