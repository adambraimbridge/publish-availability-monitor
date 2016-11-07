package feeds

import (
	"log"
	"net/url"
	"os"
	"strings"
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  - ", logPattern)
}

func NewNotificationsFeed(name string, baseUrl *url.URL, sinceDate string, expiry int, interval int, username string, password string) Feed {
	if isNotificationsPullFeed(name) {
		return &NotificationsPullFeed{
			feedName:      name,
			baseUrl:       baseUrl.String(),
			sinceDate:     sinceDate,
			expiry:        expiry + 2*interval,
			interval:      interval,
			username:      username,
			password:      password,
			notifications: make(map[string][]*Notification)}
	} else if isNotificationsPushFeed(name) {
		return &NotificationsPushFeed{
			feedName:      name,
			baseUrl:       baseUrl.String(),
			username:      username,
			password:      password,
			notifications: make(map[string][]*Notification),
			expiry:        expiry + 2*interval,
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
