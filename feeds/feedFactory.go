package feeds

import (
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

func init() {
	infoLogger = log.New()
	infoLogger.Out = os.Stdout

	errorLogger = log.New()
}

func NewNotificationsFeed(name string, baseUrl url.URL, expiry int, interval int, username string, password string) Feed {
	if isNotificationsPullFeed(name) {
		return newNotificationsPullFeed(name, baseUrl, expiry, interval, username, password)
	} else if isNotificationsPushFeed(name) {
		return newNotificationsPushFeed(name, baseUrl, expiry, interval, username, password)
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

func newNotificationsPullFeed(name string, baseUrl url.URL, expiry int, interval int, username string, password string) *NotificationsPullFeed {
	feedUrl := baseUrl.String()

	bootstrapValues := baseUrl.Query()
	bootstrapValues.Add("since", time.Now().Format(time.RFC3339))
	baseUrl.RawQuery = ""

	infoLogger.Printf("constructing NotificationsPullFeed for [%s], baseUrl = [%s], bootstrapValues = [%s]", feedUrl, baseUrl.String(), bootstrapValues.Encode())
	return &NotificationsPullFeed{
		baseNotificationsFeed{
			name,
			nil,
			feedUrl,
			username,
			password,
			expiry + 2*interval,
			make(map[string][]*Notification),
			&sync.RWMutex{},
		},
		baseUrl.String(),
		bootstrapValues.Encode(),
		&sync.Mutex{},
		interval,
		nil,
		nil,
	}
}

func newNotificationsPushFeed(name string, baseUrl url.URL, expiry int, interval int, username string, password string) *NotificationsPushFeed {
	infoLogger.Printf("constructing NotificationsPushFeed, bootstrapUrl = [%s]", baseUrl.String())
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
		false,
	}
}
