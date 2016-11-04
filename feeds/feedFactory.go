package feeds

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Financial-Times/publish-availability-monitor/checks"
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  - ", logPattern)
}

func NewNotificationsFeed(name string, httpCaller checks.HttpCaller, baseUrl *url.URL, sinceDate string, expiry int, interval int, username string, password string) Feed {
	if isNotificationsPullFeed(name) {
		return &NotificationsPullFeed{httpCaller: httpCaller,
			feedName:      name,
			baseUrl:       baseUrl.String(),
			sinceDate:     sinceDate,
			expiry:        expiry + 2*interval,
			interval:      interval,
			username:      username,
			password:      password,
			notifications: make(map[string][]*Notification)}
	} else if isNotificationsPushFeed(name) {
		return &NotificationsPushFeed{httpCaller: httpCaller,
			feedName:      name,
			baseUrl:       "http://localhost:9090/content/notifications-push",
			username:      "",
			password:      "",
			notifications: make(map[string][]*Notification)}
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
