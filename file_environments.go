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
	var validationCredentialsChanged bool
	var err error
	if validationCredentialsChanged, err = isFileChanged(validationCredentialsFileName); err != nil {
		return fmt.Errorf("Could not detect if envs file [%s] was changed. Problem was: %s", validationCredentialsFileName, err)
	}

	if !validationCredentialsChanged {
		return nil
	}

	err = updateValidationCredentials(validationCredentialsFileName)
	if err != nil {
		return fmt.Errorf("Cannot update envs. Error was: %s", err)
	}

	return nil
}

func updateEnvsIfChanged(envsFileName string, envCredentialsFileName string) error {
	var envsFileChanged, envCredentialsChanged bool
	var err error
	if envsFileChanged, err = isFileChanged(envsFileName); err != nil {
		return fmt.Errorf("Could not detect if envs file [%s] was changed. Problem was: %s", envsFileName, err)
	}

	if envCredentialsChanged, err = isFileChanged(envCredentialsFileName); err != nil {
		return fmt.Errorf("Could not detect if credentials file [%s] was changed. Problem was: %s", envCredentialsFileName, err)
	}

	if !envsFileChanged && !envCredentialsChanged {
		return nil
	}

	err = updateEnvs(envsFileName, envCredentialsFileName)
	if err != nil {
		return fmt.Errorf("Cannot update envs. Error was: %s", err)
	}

	return nil
}

func isFileChanged(fileName string) (bool, error) {
	currentHash, err := computeMD5Hash(fileName)
	if err != nil {
		return false, fmt.Errorf("Could not compute hash value for file %s. Problem was: %s", fileName, err)
	}

	var previousHash string
	var found bool
	if previousHash, found = configFilesHashValues[fileName]; !found {
		return true, nil
	}

	return previousHash != currentHash, nil
}

func computeMD5Hash(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("Could not open file with name %s. Problem was: %s", fileName, err)
	}

	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("Could not copy file with name %s to compute hash value. Problem was: %s", fileName, err)
	}

	hashValue := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashValue), nil
}

func updateEnvs(envsFileName string, envCredentialsFileName string) error {
	log.Infof("Env config files changed. Updating envs")

	envsFromFile, err := readEnvs(envsFileName)
	if err != nil {
		return fmt.Errorf("Cannot parse environments. Error was: %s", err)
	}

	validEnvs := filterInvalidEnvs(envsFromFile)

	envCredentials, err := readEnvCredentials(envCredentialsFileName)
	if err != nil {
		return fmt.Errorf("Cannot parse environments. Error was: %s", err)
	}

	removedEnvs := parseEnvsIntoMap(validEnvs, envCredentials)
	configureFileFeeds(removedEnvs)

	return nil
}

func closeFileAndUpdateHashValue(file *os.File) {
	if file == nil {
		return
	}

	fileName := file.Name()
	file.Close()
	hashValue, err := computeMD5Hash(fileName)

	if err != nil {
		log.Warn("Could not compute MD5 hash value for file %s. Problem was: %s", fileName, err)
	}
	configFilesHashValues[fileName] = hashValue
}

func updateValidationCredentials(validationCredsFileName string) error {
	log.Info("Credentials file changed. Updating validation credentials")
	credsFile, err := os.Open(validationCredsFileName)
	defer closeFileAndUpdateHashValue(credsFile)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(credsFile)
	credentials := Credentials{}
	err = jsonParser.Decode(&credentials)
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

func readEnvs(fileName string) ([]Environment, error) {
	envsFile, err := os.Open(fileName)
	defer closeFileAndUpdateHashValue(envsFile)
	if err != nil {
		return []Environment{}, err
	}

	jsonParser := json.NewDecoder(envsFile)
	envs := []Environment{}
	err = jsonParser.Decode(&envs)
	return envs, err
}

func readEnvCredentials(fileName string) ([]Credentials, error) {
	envCredsFile, err := os.Open(fileName)
	defer closeFileAndUpdateHashValue(envCredsFile)
	if err != nil {
		return []Credentials{}, err
	}

	jsonParser := json.NewDecoder(envCredsFile)
	credentials := []Credentials{}
	err = jsonParser.Decode(&credentials)

	return credentials, err
}

func isEnvInSlice(envName string, envs []Environment) bool {
	for _, env := range envs {
		if env.Name == envName {
			return true
		}
	}

	return false
}
