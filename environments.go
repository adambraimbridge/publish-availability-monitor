package main

import (
	"net/url"
	"strings"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

var (
	k8sClient kubernetes.Interface
)

func getEnvCredentials(validatorCredK8sSecretName string, credKey string) (string, error) {
	k8sSecret, err := k8sClient.CoreV1().Secrets("default").Get(validatorCredK8sSecretName)

	if err != nil {
		return "", fmt.Errorf("Error retriving validator credentials k8s secret with name %s. Error was: %s", validatorCredK8sSecretName, err.Error())
	}

	secretMap := k8sSecret.Data
	if credentials, found := secretMap[credKey]; found {
		return string(credentials[:]), nil
	}

	return "", fmt.Errorf("Entry with key %s was not found in k8s secret with name %s", credKey, validatorCredK8sSecretName)
}

func watchEnvironments(configMapName string, validatorCredK8sSecretName string, readEnvKey string, s3EnvKey string, credKey string, environments map[string]Environment) {
	fieldSelector := fmt.Sprintf("metadata.name=%s", configMapName)
	watcher, err := k8sClient.CoreV1().ConfigMaps("default").Watch(v1.ListOptions{FieldSelector: fieldSelector})

	if err != nil {
		errorLogger.Printf("Error while starting to watch envs configMap with field selector: %s. Error was: %s", fieldSelector, err.Error())
	}

	infoLogger.Print("Started watching envs configMap")
	resultChannel := watcher.ResultChan()
	for msg := range resultChannel {
		switch msg.Type {
		case watch.Added, watch.Modified:
			infoLogger.Printf("ConfigMap with name %s has been updated.", configMapName)
			k8sConfigMap := msg.Object.(*v1.ConfigMap)
			configMapData := k8sConfigMap.Data
			envReadEndpoints, found := configMapData[readEnvKey]
			if !found {
				errorLogger.Printf("Entry with key %s was not found in envs configMap. Skipping the current update on envs configMap", readEnvKey)
				continue
			}

			envS3Endpoints, found := configMapData[s3EnvKey]
			if !found {
				errorLogger.Printf("Entry with key %s was not found in envs configMap.", s3EnvKey)
			}

			envCredentials, err := getEnvCredentials(validatorCredK8sSecretName, credKey)
			if err != nil {
				errorLogger.Printf("Cannot retrieve envs credentials. Error was: %s", err.Error())
			}

			removedEnvs := parseEnvironmentsIntoMap(envReadEndpoints, envCredentials, envS3Endpoints, environments)
			configureFeeds(removedEnvs)

		case watch.Deleted:
			errorLogger.Print("Envs configMap has been removed.")
		default:
			errorLogger.Print("Error received on watch envs configMap. Channel may be full")
		}
	}

	infoLogger.Print("Env configMap watching terminated. Reconnecting...")
	watchEnvironments(configMapName, validatorCredK8sSecretName, readEnvKey, s3EnvKey, credKey, environments)
}

func updateEnvCredentials(readCredentials string) {
	envCredentials := strings.Split(readCredentials, ",")
	for _, cred := range envCredentials {
		nameAndCredentials := strings.Split(cred, ":")
		if len(nameAndCredentials) != 3 {
			warnLogger.Printf("Cannot parse credentials string: %s", cred)
			continue
		}

		envName := nameAndCredentials[0]
		if env, found := environments[envName]; found {
			env.Username = nameAndCredentials[1]
			env.Password = nameAndCredentials[2]
			environments[envName] = env
			infoLogger.Printf("Updated credentials for env with name %s", envName)
		}
	}
}

func watchCredentials(credSecretName string, validatorCredKey string, envCredentialsSecretKey string) {
	fieldSelector := fmt.Sprintf("metadata.name=%s", credSecretName)
	watcher, err := k8sClient.CoreV1().Secrets("default").Watch(v1.ListOptions{FieldSelector: fieldSelector})

	if err != nil {
		errorLogger.Printf("Error while starting to watch validatorCreds secretsMap with field selector: %s. Error was: %s", fieldSelector, err.Error())
	}

	infoLogger.Print("Started watching validator credentials secretsMap")
	resultChannel := watcher.ResultChan()
	for msg := range resultChannel {
		switch msg.Type {
		case watch.Added, watch.Modified:
			infoLogger.Printf("Secret map with name %s has been updated.", credSecretName)
			k8sSecret := msg.Object.(*v1.Secret)
			secretMap := k8sSecret.Data
			if validatorCreds, found := secretMap[validatorCredKey]; found {
				validatorCredentials = string(validatorCreds[:])
			} else {
				errorLogger.Printf("Cannot find validator credentials in secretsMap. The key to be searched is %s", validatorCredKey)
			}

			if envCreds, found := secretMap[envCredentialsSecretKey]; found {
				envCredentials := string(envCreds[:])
				updateEnvCredentials(envCredentials)
			} else {
				errorLogger.Printf("Cannot find env credentials in secretsMap. The key to be searched is %s", envCredentialsSecretKey)
			}
		case watch.Deleted:
			errorLogger.Printf("Secret map with name %s has been removed.", credSecretName)
		default:
			errorLogger.Print("Error received on watch validatorCreds secretsMap. Channel may be full")
		}
	}

	infoLogger.Print("ValidatorCreds secretsMap watching terminated. Reconnecting...")
	watchCredentials(credSecretName, validatorCredKey, envCredentialsSecretKey)
}

func DiscoverEnvironmentsAndValidators(envConfigMapName *string, credentialsSecretName *string, readEnvConfigMapKey *string, envCredentialsSecretKey *string, s3EnvConfigMapKey *string, validatorCredSecretMapKey *string, environments map[string]Environment) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to create in cluster config for k8s client. Error was: %s", err.Error()))
	}
	// creates the clientset
	k8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create k8s client, error was: %s", err.Error()))
	}

	go watchEnvironments(*envConfigMapName, *credentialsSecretName, *readEnvConfigMapKey, *s3EnvConfigMapKey, *envCredentialsSecretKey, environments)
	go watchCredentials(*credentialsSecretName, *validatorCredSecretMapKey, * envCredentialsSecretKey)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	_,err = os.Create("/etc/pam/envs/t.txt")
	if err != nil {
		panic("Cannot create file.")
	}
	go watchTestFiles(watcher)

	err = watcher.Add("/etc/pam/envs/t.txt")
	if err != nil {
		log.Fatal(err)
	}

	//err = watcher.Add("/etc/pam/credentials/read-credentials.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func watchTestFiles(watcher *fsnotify.Watcher) {
	infoLogger.Print("Starting to watch json files.")
	for {
		select {
		case event := <-watcher.Events:
			infoLogger.Printf("event: %s", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				infoLogger.Printf("modified file: %s", event.Name)
			}
		case err := <-watcher.Errors:
			infoLogger.Printf("error: %v", err)
		}
	}
}

func parseEnvironmentsIntoMap(readEnv string, readCredentials string, s3Env string, environments map[string]Environment) []string {
	envReadEndpoints := strings.Split(readEnv, ",")
	envCredentials := strings.Split(readCredentials, ",")
	envS3Endpoints := strings.Split(s3Env, ",")

	seen := make(map[string]struct{})
	for _, env := range envReadEndpoints {
		nameAndUrl := strings.SplitN(env, ":", 2)
		if len(nameAndUrl) != 2 {
			warnLogger.Printf("read-urls contain an invalid value: %s", env)
			continue
		}

		name := nameAndUrl[0]
		readUrl := nameAndUrl[1]
		seen[name] = struct{}{}

		var username string
		var password string
		for _, cred := range envCredentials {
			if strings.HasPrefix(cred, name + ":") {
				nameAndCredentials := strings.Split(cred, ":")
				username = nameAndCredentials[1]
				password = nameAndCredentials[2]
				break
			}
		}
		infoLogger.Printf("adding environment to monitoring: %v", name)
		if username == "" || password == "" {
			infoLogger.Printf("no credentials supplied for access to environment %v", name)
		}

		var s3Url string
		for _, endpoint := range envS3Endpoints {
			if strings.HasPrefix(endpoint, name + ":") {
				s3Url = strings.TrimPrefix(endpoint, name + ":")
				break
			}
		}
		if s3Url == "" {
			infoLogger.Printf("No S3 url supplied for access to environment %v", name)
		}

		environments[name] = Environment{name, readUrl, s3Url, username, password}
	}

	// now remove unseen environments
	toDelete := make([]string, 0)
	for name, _ := range environments {
		if _, exists := seen[name]; !exists {
			toDelete = append(toDelete, name)
		}
	}
	for _, name := range toDelete {
		infoLogger.Printf("removing environment from monitoring: %v", name)
		delete(environments, name)
	}

	return toDelete
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
