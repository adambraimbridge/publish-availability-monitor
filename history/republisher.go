package history

import (
	"net/http"
	"net/url"
	"strings"
)

type Republisher interface {
	Republish(uuid string, username string, password string) error
}

type jenkinsRepublisher struct {
	client               *http.Client
	republishJobBuildUrl string
	envName              string
	txPrefix             string
}

func NewJenkinsRepublisher(republishJobBuildUrl string, envName string, txPrefix string) Republisher {
	return &jenkinsRepublisher{
		&http.Client{},
		republishJobBuildUrl,
		envName,
		txPrefix}
}

func (j *jenkinsRepublisher) Republish(uuid string, username string, password string) error {
	infoLogger.Printf("republishing %s", uuid)

	form := url.Values{}
	form.Add("SOURCE_ENV", j.envName)
	form.Add("TARGET_ENV", j.envName)
	form.Add("uuidList", uuid)
	form.Add("transactionIdPrefix", j.txPrefix)

	req, err := http.NewRequest("POST", j.republishJobBuildUrl, strings.NewReader(form.Encode()))
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "UPP Publish Availability Monitor")

	resp, err := j.client.Do(req)
	if err != nil {
		infoLogger.Printf("error calling Jenkins job: %v", err)
		return err
	} else {
		infoLogger.Printf("republish job returned status %v", resp.StatusCode)
	}

	return nil
}
