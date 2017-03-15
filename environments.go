package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"

	"github.com/Financial-Times/publish-availability-monitor/feeds"
)

var (
	etcdKeysAPI  etcd.KeysAPI
	readEnvKey   *string
	s3EnvKey     *string
	credKey      *string
	validatorKey *string
)

func DiscoverEnvironmentsAndValidators(etcdPeers *string, etcdReadEnvKey *string, etcdCredKey *string, etcdS3EnvKey *string, etcdValidatorCredKey *string, environments map[string]Environment) error {
	readEnvKey = etcdReadEnvKey
	s3EnvKey = etcdS3EnvKey
	credKey = etcdCredKey
	validatorKey = etcdValidatorCredKey

	transport := &http.Transport{
		Dial: proxy.Direct.Dial,
		ResponseHeaderTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   100,
	}
	etcdCfg := etcd.Config{
		Endpoints:               strings.Split(*etcdPeers, ","),
		Transport:               transport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	}
	etcdClient, err := etcd.New(etcdCfg)
	if err != nil {
		errorLogger.Printf("Cannot load etcd configuration: [%v]", err)
		return err
	}

	etcdKeysAPI = etcd.NewKeysAPI(etcdClient)

	for len(environments) == 0 {
		if err = redefineEnvironments(environments); err != nil {
			infoLogger.Print("retry in 60s...")
			time.Sleep(time.Minute)
		}
	}

	fn := func() {
		redefineEnvironments(environments)
	}
	go watch(readEnvKey, fn)
	go watch(s3EnvKey, fn)
	go watch(credKey, fn)

	validatorCredentials = redefineValidatorCredentials()
	go watch(validatorKey, func() {
		validatorCredentials = redefineValidatorCredentials()
	})

	return nil
}

func redefineEnvironments(environments map[string]Environment) error {
	etcdReadEnvResp, err := etcdKeysAPI.Get(context.Background(), *readEnvKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *readEnvKey, err.Error())
		return err
	}

	etcdCredResp, err := etcdKeysAPI.Get(context.Background(), *credKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *credKey, err.Error())
		return err
	}

	etcdS3EnvResp, err := etcdKeysAPI.Get(context.Background(), *s3EnvKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *s3EnvKey, err.Error())
		return err
	}
	removedEnvs := parseEnvironmentsIntoMap(etcdReadEnvResp.Node.Value, etcdCredResp.Node.Value, etcdS3EnvResp.Node.Value, environments)

	configureFeeds(removedEnvs)

	return nil
}

func parseEnvironmentsIntoMap(etcdReadEnv string, etcdCred string, etcdS3Env string, environments map[string]Environment) []string {
	envReadEndpoints := strings.Split(etcdReadEnv, ",")
	envCredentials := strings.Split(etcdCred, ",")
	envS3Endpoints := strings.Split(etcdS3Env, ",")

	seen := make(map[string]struct{})
	for _, env := range envReadEndpoints {
		nameAndUrl := strings.SplitN(env, ":", 2)
		if len(nameAndUrl) != 2 {
			warnLogger.Printf("etcd read-urls contain an invalid value")
			continue
		}

		name := nameAndUrl[0]
		readUrl := nameAndUrl[1]
		seen[name] = struct{}{}

		var username string
		var password string
		for _, cred := range envCredentials {
			if strings.HasPrefix(cred, name+":") {
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
			if strings.HasPrefix(endpoint, name+":") {
				s3Url = strings.TrimPrefix(endpoint, name+":")
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

func redefineValidatorCredentials() string {
	etcdCredResp, err := etcdKeysAPI.Get(context.Background(), *validatorKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *validatorKey, err.Error())
		return ""
	}

	return etcdCredResp.Node.Value
}

func watch(etcdKey *string, fn func()) {
	watcher := etcdKeysAPI.Watcher(*etcdKey, &etcd.WatcherOptions{AfterIndex: 0, Recursive: true})
	limiter := NewEventLimiter(fn)

	for {
		_, err := watcher.Next(context.Background())
		if err != nil {
			errorLogger.Printf("Error waiting for change under %v in etcd. %v\n Sleeping 10s...", *etcdKey, err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		limiter.trigger <- true
	}
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
