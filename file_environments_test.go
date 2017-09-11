package main

import (
	"bufio"
	"fmt"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	validEnvConfig = `
		[
			{
				"name":"test-env",
				"read-url": "https://test-env.ft.com",
				"s3-url": "http://test.s3.amazonaws.com"
			}
		]`
	validEnvCredentialsConfig = `
		[
			{
				"env-name": "test-env",
				"username": "test-user",
				"password": "test-pwd"
			}
		]`
	validValidationCredentialsConfig = `
		{
			"username": "test-user",
			"password": "test-pwd"
		}`
	invalidJsonConfig = `invalid-config`
)

func TestParseEnvsIntoMap(t *testing.T) {
	envsToBeParsed := getValidEnvs()
	credentials := getValidCredentials()
	environments = make(map[string]Environment)

	removedEnvs := parseEnvsIntoMap(envsToBeParsed, credentials)

	assert.Equal(t, 0, len(removedEnvs))
	assert.Equal(t, len(envsToBeParsed), len(environments))
	envName := envsToBeParsed[1].Name
	assert.Equal(t, envName, environments[envName].Name)
	assert.Equal(t, credentials[1].Username, environments[envName].Username)
}

func TestParseEnvsIntoMapWithRemovedEnv(t *testing.T) {
	envsToBeParsed := getValidEnvs()
	credentials := getValidCredentials()
	environments = make(map[string]Environment)
	environments["removed-env"] = Environment{}

	removedEnvs := parseEnvsIntoMap(envsToBeParsed, credentials)

	assert.Equal(t, 1, len(removedEnvs))
	assert.Equal(t, len(envsToBeParsed), len(environments))
	envName := envsToBeParsed[1].Name
	assert.Equal(t, envName, environments[envName].Name)
	assert.Equal(t, credentials[1].Username, environments[envName].Username)
}

func TestParseEnvsIntoMapWithExistingEnv(t *testing.T) {
	envsToBeParsed := getValidEnvs()
	credentials := getValidCredentials()
	environments = make(map[string]Environment)
	existingEnv := envsToBeParsed[0]
	environments[existingEnv.Name] = existingEnv

	removedEnvs := parseEnvsIntoMap(envsToBeParsed, credentials)

	assert.Equal(t, 0, len(removedEnvs))
	assert.Equal(t, len(envsToBeParsed), len(environments))
	envName := envsToBeParsed[1].Name
	assert.Equal(t, envName, environments[envName].Name)
	assert.Equal(t, credentials[1].Username, environments[envName].Username)
}

func TestParseEnvsIntoMapWithNoCredentials(t *testing.T) {
	envsToBeParsed := getValidEnvs()
	credentials := []Credentials{}
	environments = make(map[string]Environment)

	removedEnvs := parseEnvsIntoMap(envsToBeParsed, credentials)

	assert.Equal(t, 0, len(removedEnvs))
	assert.Equal(t, len(envsToBeParsed), len(environments))
	envName := envsToBeParsed[1].Name
	assert.Equal(t, envName, environments[envName].Name)
}

func TestFilterInvalidEnvs(t *testing.T) {
	envsToBeFiltered := getValidEnvs()

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, len(envsToBeFiltered), len(filteredEnvs))
}

func TestFilterInvalidEnvsWithEmptyName(t *testing.T) {
	envsToBeFiltered := []Environment{
		{
			Name:     "",
			ReadUrl:  "test",
			S3Url:    "test",
			Username: "dummy",
			Password: "dummy",
		},
	}

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, 0, len(filteredEnvs))
}

func TestFilterInvalidEnvsWithEmptyReadUrl(t *testing.T) {
	envsToBeFiltered := []Environment{
		{
			Name:     "test",
			ReadUrl:  "",
			S3Url:    "test",
			Username: "dummy",
			Password: "dummy",
		},
	}

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, 0, len(filteredEnvs))
}

func TestFilterInvalidEnvsWithEmptyS3Url(t *testing.T) {
	envsToBeFiltered := []Environment{
		{
			Name:     "test",
			ReadUrl:  "test",
			S3Url:    "",
			Username: "dummy",
			Password: "dummy",
		},
	}

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, 1, len(filteredEnvs))
}

func TestFilterInvalidEnvsWithEmptyUsernameUrl(t *testing.T) {
	envsToBeFiltered := []Environment{
		{
			Name:     "test",
			ReadUrl:  "test",
			S3Url:    "test",
			Username: "",
			Password: "dummy",
		},
	}

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, 1, len(filteredEnvs))
}

func TestFilterInvalidEnvsWithEmptyPwd(t *testing.T) {
	envsToBeFiltered := []Environment{
		{
			Name:     "test",
			ReadUrl:  "test",
			S3Url:    "test",
			Username: "test",
			Password: "",
		},
	}

	filteredEnvs := filterInvalidEnvs(envsToBeFiltered)

	assert.Equal(t, 1, len(filteredEnvs))
}

func TestUpdateValidationCredentialsHappyFlow(t *testing.T) {
	fileName := prepareFile(validValidationCredentialsConfig)
	fileContents, _ := ioutil.ReadFile(fileName)
	err := updateValidationCredentials(fileContents)

	assert.Nil(t, err)
	assert.Equal(t, "test-user:test-pwd", validatorCredentials)
	os.Remove(fileName)
}

func TestUpdateValidationCredentialNilFile(t *testing.T) {
	validatorCredentials := Credentials{
		Username: "test-username",
		Password: "test-password",
	}
	err := updateValidationCredentials(nil)

	assert.NotNil(t, err)
	//make sure validationCredentials didn't change after failing call to updateValidationCredentials().
	assert.Equal(t, "test-username", validatorCredentials.Username)
	assert.Equal(t, "test-password", validatorCredentials.Password)
}

func TestUpdateValidationCredentialsInvalidConfig(t *testing.T) {
	fileName := prepareFile(invalidJsonConfig)
	validatorCredentials := Credentials{
		Username: "test-username",
		Password: "test-password",
	}
	fileContents, _ := ioutil.ReadFile(fileName)
	err := updateValidationCredentials(fileContents)
	assert.NotNil(t, err)
	//make sure validationCredentials didn't change after failing call to updateValidationCredentials().
	assert.Equal(t, "test-username", validatorCredentials.Username)
	assert.Equal(t, "test-password", validatorCredentials.Password)
	os.Remove(fileName)
}

func TestConfigureFeedsWithEmptyListOfMetrics(t *testing.T) {
	subscribedFeeds["test-feed"] = []feeds.Feed{
		MockFeed{},
	}
	appConfig = &AppConfig{}

	configureFileFeeds([]string{"test-feed"})

	assert.Equal(t, 0, len(subscribedFeeds))
}

func TestUpdateEnvsHappyFlow(t *testing.T) {
	subscribedFeeds["test-feed"] = []feeds.Feed{
		MockFeed{},
	}
	appConfig = &AppConfig{}
	envsFileName := prepareFile(validEnvConfig)
	envsFileContents, _ := ioutil.ReadFile(envsFileName)

	envCredsFileName := prepareFile(validEnvCredentialsConfig)
	credsFileContents, _ := ioutil.ReadFile(envCredsFileName)

	err := updateEnvs(envsFileContents, credsFileContents)

	assert.Nil(t, err)
	os.Remove(envsFileName)
	os.Remove(envCredsFileName)
}

func TestUpdateEnvsHappyNilEnvsFile(t *testing.T) {
	envCredsFileName := prepareFile(validEnvCredentialsConfig)
	credsFileContents, _ := ioutil.ReadFile(envCredsFileName)
	err := updateEnvs(nil, credsFileContents)

	assert.NotNil(t, err)
	os.Remove(envCredsFileName)
}

func TestUpdateEnvsNilEnvCredentialsFile(t *testing.T) {
	envsFileName := prepareFile(validEnvConfig)
	envsFileContents, _ := ioutil.ReadFile(envsFileName)

	err := updateEnvs(envsFileContents, nil)

	assert.NotNil(t, err)
	os.Remove(envsFileName)
}

func prepareFile(fileContent string) string {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		panic("Cannot create temp file.")
	}

	writer := bufio.NewWriter(file)
	defer file.Close()
	fmt.Fprintln(writer, fileContent)
	writer.Flush()
	return file.Name()
}

func getValidEnvs() []Environment {
	return []Environment{
		{
			Name:    "test",
			ReadUrl: "test-url",
			S3Url:   "test-s3-url",
		},
		{
			Name:    "test2",
			ReadUrl: "test-url2",
			S3Url:   "test-s3-url2",
		},
	}
}

func getValidCredentials() []Credentials {
	return []Credentials{
		{
			EnvName:  "test",
			Username: "dummy-user",
			Password: "dummy-pwd",
		},
		{
			EnvName:  "test2",
			Username: "dummy-user2",
			Password: "dummy-pwd2",
		},
	}
}

type MockFeed struct{}

func (f MockFeed) Start() {}
func (f MockFeed) Stop()  {}
func (f MockFeed) FeedName() string {
	return ""
}
func (f MockFeed) FeedURL() string {
	return ""
}
func (f MockFeed) FeedType() string {
	return ""
}
func (f MockFeed) SetCredentials(username string, password string) {}
func (f MockFeed) NotificationsFor(uuid string) []*feeds.Notification {
	return nil
}