package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
)

var (
	etcdKeysAPI  etcd.KeysAPI
	envKey       *string
	credKey      *string
	validatorKey *string
)

func DiscoverEnvironmentsAndValidators(etcdPeers *string, etcdEnvKey *string, etcdCredKey *string, etcdValidatorCredKey *string, environments map[string]Environment) error {
	envKey = etcdEnvKey
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
	go watch(envKey, fn)
	go watch(credKey, fn)

	validatorCredentials = redefineValidatorCredentials()
	go watch(validatorKey, func() {
		validatorCredentials = redefineValidatorCredentials()
	})

	return nil
}

func redefineEnvironments(environments map[string]Environment) error {
	etcdEnvResp, err := etcdKeysAPI.Get(context.Background(), *envKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *envKey, err.Error())
		return err
	}

	etcdCredResp, err := etcdKeysAPI.Get(context.Background(), *credKey, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", *credKey, err.Error())
		return err
	}

	parseEnvironmentsIntoMap(etcdEnvResp.Node.Value, etcdCredResp.Node.Value, environments)

	startFeeds()

	return nil
}

func parseEnvironmentsIntoMap(etcdEnv string, etcdCred string, environments map[string]Environment) {
	envReadEndpoints := strings.Split(etcdEnv, ",")
	envCredentials := strings.Split(etcdCred, ",")

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

		environments[name] = Environment{name, readUrl, username, password}
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

func startFeeds() {
	for _, metric := range appConfig.MetricConf {
		if metric.Alias == "notifications" {
			for _, env := range environments {
				httpCaller := checks.NewHttpCaller()
				endpointUrl, err := url.Parse(env.ReadUrl + metric.Endpoint)
				if err != nil {
					errorLogger.Printf("Cannot parse url [%v], error: [%v]", metric.Endpoint, err.Error())
					continue
				}

				sinceDate := time.Now().Format(time.RFC3339)
				infoLogger.Printf("since %v", sinceDate)
				interval := appConfig.Threshold / metric.Granularity

				f := feeds.NewNotificationsPullFeed(httpCaller, endpointUrl, sinceDate, appConfig.Threshold, interval, env.Username, env.Password)

				var envFeeds []feeds.Feed
				var found bool
				if envFeeds, found = subscribedFeeds[env.Name]; !found {
					envFeeds = make([]feeds.Feed, 0)
				}
				subscribedFeeds[env.Name] = append(envFeeds, f)

				f.Start()
			}
		}
	}
}
