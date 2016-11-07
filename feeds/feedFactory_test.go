package feeds

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPullFeed(t *testing.T) {
	baseUrl, _ := url.Parse("http://www.example.org/")

	actual := NewNotificationsFeed("notifications", baseUrl, "2016-11-01T00:00:00Z", 10, 10, "", "")
	assert.IsType(t, (*NotificationsPullFeed)(nil), actual, "expected a NotificationsPullFeed")
}

func TestNewPushFeed(t *testing.T) {
	baseUrl, _ := url.Parse("http://www.example.org/")

	actual := NewNotificationsFeed("notifications-push", baseUrl, "2016-11-01T00:00:00Z", 10, 10, "", "")
	assert.IsType(t, (*NotificationsPushFeed)(nil), actual, "expected a NotificationsPushFeed")
}
