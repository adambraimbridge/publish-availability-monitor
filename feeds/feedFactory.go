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
			baseNotificationsFeed{
				name,
				nil,
				baseUrl.String(),
				username,
				password,
				expiry + 2*interval,
				make(map[string][]*Notification),
				&sync.RWMutex{},
			},
			sinceDate,
			&sync.Mutex{},
			interval,
			nil,
			nil,
		}
	} else if isNotificationsPushFeed(name) {
		return &NotificationsPushFeed{
			baseNotificationsFeed{
				name,
				nil,
				baseUrl.String(),
				username,
				password,
				expiry + 2*interval,
				make(map[string][]*Notification),
				&sync.RWMutex{},
			},
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
