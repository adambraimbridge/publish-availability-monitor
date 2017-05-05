package main

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
	"github.com/fsnotify/fsnotify"
	"net/url"
	"os"
)

func watchConfigFiles(envsFileName string, envCredentialsFileName string, validationCredentialsFileName string) {
	var watcher *fsnotify.Watcher
	var receivedRemoveEvent bool
	var err error
	for {
		if watcher == nil {
			infoLogger.Print("Creating new watcher.")
			// Intialize the watcher if it has not been initialized yet.
			if watcher, err = fsnotify.NewWatcher(); err != nil {
				errorLogger.Printf("failed to create fsnotify watcher: %v", err)
			}

			defer watcher.Close()

			//watch configuration files
			filesToBeWatched := []string{envsFileName, envCredentialsFileName, validationCredentialsFileName}
			for _, fileToBeWatched := range filesToBeWatched {
				if err := watcher.Add(fileToBeWatched); err != nil {
					errorLogger.Printf("failed to watch file %q: %v", fileToBeWatched, err)
					continue
				}

				infoLogger.Printf("Started watching file with name %s", fileToBeWatched)
			}
		}

		// Wait until the next log change.
		receivedRemoveEvent, err = monitorFilesForChanges(watcher)
		if err != nil {
			errorLogger.Printf("Err received on wating test logs. Err is %s", err)
			return
		}

		if receivedRemoveEvent {
			errorLogger.Print("Received remove event, removing watcher.")
			watcher.Close()
			watcher = nil
			continue
		}

		infoLogger.Print("Config files have been changed.")
		err := updateEnvsAndValidationCredentials(envsFileName, envCredentialsFileName, validationCredentialsFileName)
		if err != nil {
			errorLogger.Printf("Cannot update configuration, error was: %s",err)
		}
	}
}

func updateEnvsAndValidationCredentials(envsFileName string, envCredentialsFileName string, validationCredentialsFileName string) error {
	err := updateEnvs(envsFileName, envCredentialsFileName)
	if err != nil {
		return fmt.Errorf("Cannot update envs. Error was: %s", err)
	}

	err = updateValidationCredentials(validationCredentialsFileName)
	if err != nil {
		fmt.Errorf("Cannot update envs. Error was: %s", err)
	}

	return nil
}

func monitorFilesForChanges(w *fsnotify.Watcher) (bool, error) {
	infoLogger.Print("Waiting for changes on file")
	errRetry := 5
	for {
		select {
		case e := <-w.Events:
			switch e.Op {
			case fsnotify.Write:
				return false, nil
			case fsnotify.Remove:
				errorLogger.Printf("Received remove event: %v", e)
				return true, nil
			default:
				errorLogger.Printf("Unexpected fsnotify event: %v, retrying...", e)
			}
		case err := <-w.Errors:
			errorLogger.Printf("Fsnotify watch error: %v, %d error retries remaining", err, errRetry)
			if errRetry == 0 {
				return false, err
			}
			errRetry--
		}
	}
}

func updateEnvs(envsFileName string, envCredentialsFileName string) error {
	infoLogger.Print("Updating envs")

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
	configureFeeds(removedEnvs)

	return nil
}

func updateValidationCredentials(validationCredsFileName string) error {
	infoLogger.Print("Updating validation credentials")
	data, err := os.Open(validationCredsFileName)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(data)
	credentials := Credentials{}
	err = jsonParser.Decode(&credentials)
	if err != nil {
		return err
	}

	validatorCredentials = credentials
	return nil
}

func configureFeeds(removedEnvs []string) {
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
					errorLogger.Printf("Cannot parse url [%v], error: [%v]", metric.Endpoint, err.Error())
					continue
				}

				interval := appConfig.Threshold / metric.Granularity

				if f := feeds.NewNotificationsFeed(metric.Alias, *endpointUrl, appConfig.Threshold, interval, env.Username, env.Password); f != nil {
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
			errorLogger.Printf("Env %v has an empty name, skipping it", env)
			continue
		}

		//envs without read-url are invalid
		if env.ReadUrl == "" {
			errorLogger.Printf("Env with name %s does not have readUrl, skipping it", env.Name)
			continue
		}

		//envs without s3 are still valid, but still a heads up is given.
		if env.S3Url == "" {
			infoLogger.Printf("Env with name %s does not have s3 url.", env.S3Url)
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
			infoLogger.Printf("No credentials provided for env with name %s", env.Name)
		}
	}

	//remove envs that don't exist anymore
	removedEnvs := make([]string, 0)
	for envName := range environments {
		if !isEnvInSlice(envName, envs) {
			infoLogger.Printf("removing environment from monitoring: %v", envName)
			delete(environments, envName)
			removedEnvs = append(removedEnvs, envName)
		}
	}

	//update envs
	for _, env := range envs {
		envName := env.Name
		environments[envName] = env
		infoLogger.Printf("Added environment to monitoring: %s", envName)
	}

	return removedEnvs
}

func readEnvs(fileName string) ([]Environment, error) {
	data, err := os.Open(fileName)
	if err != nil {
		return []Environment{}, err
	}

	jsonParser := json.NewDecoder(data)
	envs := []Environment{}
	err = jsonParser.Decode(&envs)
	return envs, err
}

func readEnvCredentials(fileName string) ([]Credentials, error) {
	data, err := os.Open(fileName)
	if err != nil {
		return []Credentials{}, err
	}

	jsonParser := json.NewDecoder(data)
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
