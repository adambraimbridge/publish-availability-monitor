package main

import (
	"net/http"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
)

var (
	etcdKeysAPI etcd.KeysAPI
	envKey      *string
	credKey     *string
)

func DiscoverEnvironments(etcdPeers *string, etcdEnvKey *string, etcdCredKey *string, environments map[string]Environment) error {
	envKey = etcdEnvKey
	credKey = etcdCredKey

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

	go watch(envKey, environments)
	go watch(credKey, environments)

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

func watch(etcdKey *string, environments map[string]Environment) {
	watcher := etcdKeysAPI.Watcher(*etcdKey, &etcd.WatcherOptions{AfterIndex: 0, Recursive: true})
	limiter := NewEventLimiter(func() {
		redefineEnvironments(environments)
	})

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
