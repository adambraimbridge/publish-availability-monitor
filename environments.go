package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
)

const (
	environmentsKeyPre = "/ft/config/publish-availability-monitor/delivery-environments"
	readUrlSuffix      = "/read_url"
	credentialsKeyPre  = "/ft/_credentials/coco-delivery"
	credUsernameKey    = credentialsKeyPre + "/%s/username"
	credPasswordKey    = credentialsKeyPre + "/%s/password"
)

var (
	etcdKeysAPI etcd.KeysAPI
)

func DiscoverEnvironments(etcdPeers *string) (map[string]Environment, error) {
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
		return nil, err
	}

	etcdKeysAPI = etcd.NewKeysAPI(etcdClient)

	environments = make(map[string]Environment)

	err = redefineEnvironments(environments)
	if err != nil {
		return nil, err
	}

	go watchEnvironments(environments)
	go watchCredentials(environments)

	return environments, nil
}

func redefineEnvironments(environments map[string]Environment) error {
	etcdResp, err := etcdKeysAPI.Get(context.Background(), environmentsKeyPre, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", environmentsKeyPre, err.Error())
		return err
	}
	if !etcdResp.Node.Dir {
		errorLogger.Printf("[%v] is not a directory", etcdResp.Node.Key)
		return err
	}

	seen := make(map[string]struct{})
	for _, envNode := range etcdResp.Node.Nodes {
		if !envNode.Dir {
			warnLogger.Printf("[%v] is not a directory", envNode.Key)
			continue
		}
		name := filepath.Base(envNode.Key)
		seen[name] = struct{}{}
		pathResp, err := etcdKeysAPI.Get(context.Background(), envNode.Key+readUrlSuffix, &etcd.GetOptions{Sort: true})
		if err != nil {
			warnLogger.Printf("Failed to get read url path from %v: %v.", envNode.Key, err.Error())
			return err
		}
		readUrl := pathResp.Node.Value

		var username string
		var password string
		credUserResp, err := etcdKeysAPI.Get(context.Background(), fmt.Sprintf(credUsernameKey, name), nil)
		if err == nil {
			credPassResp, err := etcdKeysAPI.Get(context.Background(), fmt.Sprintf(credPasswordKey, name), nil)

			if err == nil {
				username = credUserResp.Node.Value
				password = credPassResp.Node.Value
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

	return nil
}

func watchEnvironments(environments map[string]Environment) {
	watch(environmentsKeyPre, environments)
}

func watchCredentials(environments map[string]Environment) {
	watch(credentialsKeyPre, environments)
}

func watch(etcdDir string, environments map[string]Environment) {
	watcher := etcdKeysAPI.Watcher(etcdDir, &etcd.WatcherOptions{AfterIndex: 0, Recursive: true})
	limiter := NewEventLimiter(func() {
		redefineEnvironments(environments)
	})

	for {
		_, err := watcher.Next(context.Background())
		if err != nil {
			errorLogger.Printf("Error waiting for change under %v in etcd. %v\n Sleeping 10s...", environmentsKeyPre, err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		limiter.trigger <- true
	}
}
