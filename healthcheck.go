package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	"io/ioutil"
	"net/http"
)

type Healthcheck struct {
	client http.Client
	config AppConfig
}

func (h *Healthcheck) CheckHealth() func(w http.ResponseWriter, r *http.Request) {
	return fthealth.HandlerParallel("Dependent services healthcheck", "Checks if all the dependent services are reachable and healthy.", h.messageQueueProxyReachable())
}

func (h *Healthcheck) Gtg(writer http.ResponseWriter, req *http.Request) {
	healthChecks := []func() error{h.checkAggregateMessageQueueProxiesReachable}

	for _, hCheck := range healthChecks {
		if err := hCheck(); err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
}

func (h *Healthcheck) messageQueueProxyReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
		Name:             "MessageQueueProxyReachable",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-runbook",
		Severity:         1,
		TechnicalSummary: "Message queue proxy is not reachable/healthy",
		Checker:          h.checkAggregateMessageQueueProxiesReachable,
	}

}

func (h *Healthcheck) checkAggregateMessageQueueProxiesReachable() error {

	addresses := h.config.QueueConf.Addrs
	errMsg := ""
	for i := 0; i < len(addresses); i++ {
		error := h.checkMessageQueueProxyReachable(addresses[i])
		if error == nil {
			return nil
		} else {
			errMsg = errMsg + fmt.Sprintf("For %s there is an error %v \n", addresses[i], error.Error())
		}
	}

	return errors.New(errMsg)

}

func (h *Healthcheck) checkMessageQueueProxyReachable(address string) error {
	req, err := http.NewRequest("GET", address+"/topics", nil)
	if err != nil {
		warn.Printf("Could not connect to proxy: %v", err.Error())
		return err
	}

	if len(h.config.QueueConf.AuthorizationKey) > 0 {
		req.Header.Add("Authorization", h.config.QueueConf.AuthorizationKey)
	}

	if len(h.config.QueueConf.Queue) > 0 {
		req.Host = h.config.QueueConf.Queue
	}

	resp, err := h.client.Do(req)
	if err != nil {
		warn.Printf("Could not connect to proxy: %v", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Proxy returned status: %d", resp.StatusCode)
		return errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return checkIfTopicIsPresent(body, h.config.QueueConf.Topic)

}

func checkIfTopicIsPresent(body []byte, searchedTopic string) error {
	var topics []string

	err := json.Unmarshal(body, &topics)
	if err != nil {
		return errors.New(fmt.Sprintf("Error occured and topic could not be found. %v", err.Error()))
	}

	for _, topic := range topics {
		if topic == searchedTopic {
			return nil
		}
	}

	return errors.New("Topic was not found")
}
