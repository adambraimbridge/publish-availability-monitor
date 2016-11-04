package feeds

// ignore unused fields (e.g. type, apiUrl)
type Notification struct {
	PublishReference string
	LastModified     string
	ID               string
}

// ignore unused field (e.g. rel)
type Link struct {
	Href string
}

type Feed interface {
	Start()
	Stop()
	FeedName() string
	FeedType() string
	SetCredentials(username string, password string)
	NotificationsFor(uuid string) []*Notification
}
