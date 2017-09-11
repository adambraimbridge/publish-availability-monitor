package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/feeds"
	log "github.com/Sirupsen/logrus"
	"errors"
)

func watchConfigFiles(envsFileName string, envCredentialsFileName string, validationCredentialsFileName string, configRefreshPeriod int) {
	ticker := time.NewTicker(time.Minute * time.Duration(configRefreshPeriod))

	for range ticker.C {
		err := updateEnvsIfChanged(envsFileName, envCredentialsFileName)
		if err != nil {
			log.Errorf("Could not update envs config, error was: %s", err)
		}

		err = updateValidationCredentialsIfChanged(validationCredentialsFileName)
		if err != nil {
			log.Errorf("Could not update validation credentials config, error was: %s", err)
		}
	}
}

func updateValidationCredentialsIfChanged(validationCredentialsFileName string) error {
	credsFile, err := os.Open(validationCredentialsFileName)
	if err != nil {
		return fmt.Errorf("could not open creds file [%s] because [%s]", envsFileName, err)
	}
	defer credsFile.Close()

	var validationCredentialsChanged bool
	var credsNewHash string
	if validationCredentialsChanged, credsNewHash, err = isFileChanged(credsFile); err != nil {
		return fmt.Errorf("could not detect if creds file [%s] was changed because: [%s]", validationCredentialsFileName, err)
	}

	if !validationCredentialsChanged {
		return nil
	}

	err = updateValidationCredentials(credsFile)
	if err != nil {
		return fmt.Errorf("cannot update validation credentials because [%s]", err)
	}

	configFilesHashValues[credsFile.Name()] = credsNewHash
	return nil
}

func updateEnvsIfChanged(envsFileName string, envCredentialsFileName string) error {
	var envsFileChanged, envCredentialsChanged bool
	var envsNewHash, credsNewHash string

	envsFile, err := os.Open(envsFileName)
	if err != nil {
		return fmt.Errorf("could not open envs file [%s] because [%s]", envsFileName, err)
	}
	defer envsFile.Close()

	if envsFileChanged, envsNewHash, err = isFileChanged(envsFile); err != nil {
		return fmt.Errorf("could not detect if envs file [%s] was changed because [%s]", envsFileName, err)
	}

	credsFile, err := os.Open(envCredentialsFileName)
	if err != nil {
		return fmt.Errorf("could not open creds file [%s] because [%s]", envCredentialsFileName, err)
	}
	defer credsFile.Close()

	if envCredentialsChanged, credsNewHash, err = isFileChanged(credsFile); err != nil {
		return fmt.Errorf("could not detect if credentials file [%s] was changed because [%s]", envCredentialsFileName, err)
	}

	if !envsFileChanged && !envCredentialsChanged {
		return nil
	}

	err = updateEnvs(envsFile, credsFile)
	if err != nil {
		return fmt.Errorf("cannot update environments and credentials because [%s]", err)
	}
	configFilesHashValues[envsFile.Name()] = envsNewHash
	configFilesHashValues[credsFile.Name()] = credsNewHash
	return nil
}

func isFileChanged(file *os.File) (bool, string, error) {
	currentHash, err := computeMD5Hash(file)
	if err != nil {
		return false, "", fmt.Errorf("could not compute hash value for file [%s] because [%s]", file.Name(), err)
	}

	previousHash, found := configFilesHashValues[file.Name()]
	if found && previousHash == currentHash {
		return false, previousHash, nil
	}

	return true, currentHash, nil
}

func computeMD5Hash(file *os.File) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("could not copy file [%s] to compute hash value because [%s]", file.Name(), err)
	}
	hashValue := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashValue), nil
}

func updateEnvs(envsFile *os.File, credsFile *os.File) error {
	if envsFile == nil {
		return errors.New("cannot update envs because envs file is nil")
	}
	if credsFile == nil{
		return errors.New("cannot update env credentials because credentials file is nil")
	}
	log.Infof("Env config files changed. Updating envs")

	jsonParser := json.NewDecoder(envsFile)
	envsFromFile := []Environment{}
	err := jsonParser.Decode(&envsFromFile)
	if err != nil {
		return fmt.Errorf("cannot parse environmente because [%s]", err)
	}

	validEnvs := filterInvalidEnvs(envsFromFile)

	jsonParser = json.NewDecoder(credsFile)
	envCredentials := []Credentials{}
	err = jsonParser.Decode(&envCredentials)

	if err != nil {
		return fmt.Errorf("cannot parse credentials because [%s]", err)
	}

	removedEnvs := parseEnvsIntoMap(validEnvs, envCredentials)
	configureFileFeeds(removedEnvs)

	return nil
}

func updateValidationCredentials(credsFile *os.File) error {
	if credsFile == nil{
		return errors.New("cannot update validation credentials from nil file")
	}
	log.Info("Credentials file changed. Updating validation credentials")

	jsonParser := json.NewDecoder(credsFile)
	credentials := Credentials{}
	err := jsonParser.Decode(&credentials)
	if err != nil {
		return err
	}
	validatorCredentials = credentials.Username + ":" + credentials.Password
	return nil
}

func configureFileFeeds(removedEnvs []string) {
	for _, envName := range removedEnvs {
		feeds, found := subscribedFeeds[envName]
		if found {
			for _, f := range feeds {
				f.Stop()
			}
		}

		delete(subscribedFeeds, envName)
	}

	for _, metric := range appConfig.MetricConf {
		for _, env := range environments {
			var envFeeds []feeds.Feed
			var found bool
			if envFeeds, found = subscribedFeeds[env.Name]; !found {
				envFeeds = make([]feeds.Feed, 0)
			}

			found = false
			for _, f := range envFeeds {
				if f.FeedName() == metric.Alias {
					f.SetCredentials(env.Username, env.Password)
					found = true
					break
				}
			}

			if !found {
				endpointUrl, err := url.Parse(env.ReadUrl + metric.Endpoint)
				if err != nil {
					log.Errorf("Cannot parse url [%v], error: [%v]", metric.Endpoint, err.Error())
					continue
				}

				interval := appConfig.Threshold / metric.Granularity

				if f := feeds.NewNotificationsFeed(metric.Alias, *endpointUrl, appConfig.Threshold, interval, env.Username, env.Password, metric.ApiKey); f != nil {
					subscribedFeeds[env.Name] = append(envFeeds, f)
					f.Start()
				}
			}
		}
	}
}

func filterInvalidEnvs(envsFromFile []Environment) []Environment {
	var validEnvs []Environment
	for _, env := range envsFromFile {
		//envs without name are invalid
		if env.Name == "" {
			log.Errorf("Env %v has an empty name, skipping it", env)
			continue
		}

		//envs without read-url are invalid
		if env.ReadUrl == "" {
			log.Errorf("Env with name %s does not have readUrl, skipping it", env.Name)
			continue
		}

		//envs without s3 are still valid, but still a heads up is given.
		if env.S3Url == "" {
			log.Errorf("Env with name %s does not have s3 url.", env.S3Url)
		}

		validEnvs = append(validEnvs, env)
	}

	return validEnvs
}

func parseEnvsIntoMap(envs []Environment, envCredentials []Credentials) []string {
	//enhance envs with credentials
	for i, env := range envs {
		for _, envCredentials := range envCredentials {
			if env.Name == envCredentials.EnvName {
				envs[i].Username = envCredentials.Username
				envs[i].Password = envCredentials.Password
				break
			}
		}

		if envs[i].Username == "" || envs[i].Password == "" {
			log.Infof("No credentials provided for env with name %s", env.Name)
		}
	}

	//remove envs that don't exist anymore
	removedEnvs := make([]string, 0)
	for envName := range environments {
		if !isEnvInSlice(envName, envs) {
			log.Infof("removing environment from monitoring: %v", envName)
			delete(environments, envName)
			removedEnvs = append(removedEnvs, envName)
		}
	}

	//update envs
	for _, env := range envs {
		envName := env.Name
		environments[envName] = env
		log.Infof("Added environment to monitoring: %s", envName)
	}

	return removedEnvs
}

func isEnvInSlice(envName string, envs []Environment) bool {
	for _, env := range envs {
		if env.Name == envName {
			return true
		}
	}

	return false
}
