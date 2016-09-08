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
	credUsernameKey    = "/ft/_credentials/coco-delivery/%s/username"
	credPasswordKey    = "/ft/_credentials/coco-delivery/%s/password"
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

	etcdKeysAPI := etcd.NewKeysAPI(etcdClient)

	environments = make(map[string]Environment)
	etcdResp, err := etcdKeysAPI.Get(context.Background(), environmentsKeyPre, &etcd.GetOptions{Sort: true})
	if err != nil {
		errorLogger.Printf("Failed to get value from %v: %v.", environmentsKeyPre, err.Error())
		return nil, err
	}
	if !etcdResp.Node.Dir {
		errorLogger.Printf("[%v] is not a directory", etcdResp.Node.Key)
		return nil, err
	}
	for _, envNode := range etcdResp.Node.Nodes {
		if !envNode.Dir {
			warnLogger.Printf("[%v] is not a directory", envNode.Key)
			continue
		}
		name := filepath.Base(envNode.Key)
		pathResp, err := etcdKeysAPI.Get(context.Background(), envNode.Key+readUrlSuffix, &etcd.GetOptions{Sort: true})
		if err != nil {
			warnLogger.Printf("Failed to get read url path from %v: %v.", envNode.Key, err.Error())
			return nil, err
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

	return environments, nil
}
