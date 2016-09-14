package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	fthealth "github.com/Financial-Times/go-fthealth/v1a"
)

// Healthcheck offers methods to measure application health.
type Healthcheck struct {
	client          http.Client
	config          AppConfig
	metricContainer *publishHistory
}

type readEnvironmentHealthcheck struct {
	env    Environment
	client http.Client
}

var readCheckEndpoints = map[string]func(string) (string, error){
	"S3": buildAwsHealthcheckUrl,
	// only exceptions need to be listed here - everything else will default to standard FT healthcheck URLs
}

var noReadEnvironments = fthealth.Check{
	BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
	Name:             "ReadEnvironments",
	PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-run-book",
	Severity:         1,
	TechnicalSummary: "There are no read environments to monitor. This could be because none have been configured, or that etcd is not reachable/healthy",
	Checker: func() (string, error) {
		return "", errors.New("There are no read environments to monitor.")
	},
}

func (h *Healthcheck) checkHealth(writer http.ResponseWriter, req *http.Request) {
	checks := make([]fthealth.Check, 3)
	checks[0] = h.messageQueueProxyReachable()
	checks[1] = h.reflectPublishFailures()
	checks[2] = h.validationServicesReachable()

	readEnvironmentChecks := h.readEnvironmentsReachable()
	if len(readEnvironmentChecks) == 0 {
		checks = append(checks, noReadEnvironments)
	} else {
		for _, hc := range readEnvironmentChecks {
			checks = append(checks, hc)
		}
	}

	fthealth.HandlerParallel(
		"Dependent services healthcheck", "Checks if all the dependent services are reachable and healthy.",
		checks...,
	)(writer, req)
}

func (h *Healthcheck) gtg(writer http.ResponseWriter, req *http.Request) {
	healthChecks := []func() (string, error){h.checkAggregateMessageQueueProxiesReachable, h.checkValidationServicesReachable}

	for _, hCheck := range healthChecks {
		if _, err := hCheck(); err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}
}

func (h *Healthcheck) messageQueueProxyReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics are not recorded. This will impact the SLA measurement.",
		Name:             "MessageQueueProxyReachable",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-run-book",
		Severity:         1,
		TechnicalSummary: "Message queue proxy is not reachable/healthy",
		Checker:          h.checkAggregateMessageQueueProxiesReachable,
	}

}

func (h *Healthcheck) checkAggregateMessageQueueProxiesReachable() (string, error) {

	addresses := h.config.QueueConf.Addrs
	errMsg := ""
	for i := 0; i < len(addresses); i++ {
		msg, error := h.checkMessageQueueProxyReachable(addresses[i])
		if error == nil {
			return msg, nil
		}
		errMsg = errMsg + fmt.Sprintf("For %s there is an error %v \n", addresses[i], error.Error())
	}

	return "", errors.New(errMsg)

}

func (h *Healthcheck) checkMessageQueueProxyReachable(address string) (string, error) {
	req, err := http.NewRequest("GET", address+"/topics", nil)
	if err != nil {
		warnLogger.Printf("Could not connect to proxy: %v", err.Error())
		return "", err
	}

	if len(h.config.QueueConf.AuthorizationKey) > 0 {
		req.Header.Add("Authorization", h.config.QueueConf.AuthorizationKey)
	}

	if len(h.config.QueueConf.Queue) > 0 {
		req.Host = h.config.QueueConf.Queue
	}

	resp, err := h.client.Do(req)
	if err != nil {
		warnLogger.Printf("Could not connect to proxy: %v", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Proxy returned status: %d", resp.StatusCode)
		return "", errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return checkIfTopicIsPresent(body, h.config.QueueConf.Topic)

}

func checkIfTopicIsPresent(body []byte, searchedTopic string) (string, error) {
	var topics []string

	err := json.Unmarshal(body, &topics)
	if err != nil {
		return "", fmt.Errorf("Error occured and topic could not be found. %v", err.Error())
	}

	for _, topic := range topics {
		if topic == searchedTopic {
			return "", nil
		}
	}

	return "", errors.New("Topic was not found")
}

func (h *Healthcheck) reflectPublishFailures() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "At least two of the last 10 publishes failed. This will reflect in the SLA measurement.",
		Name:             "ReflectPublishFailures",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-run-book",
		Severity:         1,
		TechnicalSummary: "Publishes did not meet the SLA measurments",
		Checker:          h.checkForPublishFailures,
	}

}

func (h *Healthcheck) checkForPublishFailures() (string, error) {
	metricContainer.RLock()
	failures := 0
	for i := 0; i < len(metricContainer.publishMetrics); i++ {
		if !metricContainer.publishMetrics[i].publishOK {
			failures++
		}
	}
	metricContainer.RUnlock()

	failureThreshold := 2 //default
	if h.config.HealthConf.FailureThreshold != 0 {
		failureThreshold = h.config.HealthConf.FailureThreshold
	}
	if failures >= failureThreshold {
		return "", fmt.Errorf("%d publish failures happened during the last 10 publishes", failures)
	}
	return "", nil
}

func (h *Healthcheck) validationServicesReachable() fthealth.Check {
	return fthealth.Check{
		BusinessImpact:   "Publish metrics might not be correct. False positive failures might be recorded. This will impact the SLA measurement.",
		Name:             "validationServicesReachable",
		PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-run-book",
		Severity:         1,
		TechnicalSummary: "Validation services are not reachable/healthy",
		Checker:          h.checkValidationServicesReachable,
	}
}

func (h *Healthcheck) checkValidationServicesReachable() (string, error) {
	endpoints := h.config.ValidationEndpoints
	var wg sync.WaitGroup
	hcErrs := make(chan error, len(endpoints))
	for _, url := range endpoints {
		wg.Add(1)
		go checkServiceReachable("validation", url, h.client, hcErrs, &wg)
	}
	wg.Wait()
	close(hcErrs)
	for err := range hcErrs {
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func checkServiceReachable(serviceType string, serviceURL string, client http.Client, hcRes chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	var healthURL string
	var err error
	if fn, ok := readCheckEndpoints[serviceType]; ok {
		healthURL, err = fn(serviceURL)
	} else {
		healthURL, err = buildFtHealthcheckUrl(serviceURL)
	}
	if err != nil {
		hcRes <- fmt.Errorf("Service URL: [%s]. Error: [%v]", serviceURL, err)
		return
	}

	resp, err := client.Get(healthURL)
	if err != nil {
		hcRes <- fmt.Errorf("Healthcheck URL: [%s]. Error: [%v]", healthURL, err)
		return
	}
	defer cleanupResp(resp)
	if resp.StatusCode != 200 {
		hcRes <- fmt.Errorf("Not healthy statusCode received: [%d]", resp.StatusCode)
		return
	}
	hcRes <- nil
}

func (h *Healthcheck) readEnvironmentsReachable() []fthealth.Check {
	hc := make([]fthealth.Check, len(environments))

	i := 0
	for _, env := range environments {
		hc[i] = fthealth.Check{
			BusinessImpact:   "Publish metrics might not be correct. False positive failures might be recorded. This will impact the SLA measurement.",
			Name:             env.Name + " readEndpointsReachable",
			PanicGuide:       "https://sites.google.com/a/ft.com/technology/systems/dynamic-semantic-publishing/extra-publishing/publish-availability-monitor-run-book",
			Severity:         1,
			TechnicalSummary: "Read services are not reachable/healthy",
			Checker:          (&readEnvironmentHealthcheck{env, h.client}).checkReadEnvironmentReachable,
		}
		i++
	}
	return hc
}

func (h *readEnvironmentHealthcheck) checkReadEnvironmentReachable() (string, error) {
	var wg sync.WaitGroup
	hcErrs := make(chan error, len(appConfig.MetricConf))

	for _, metric := range appConfig.MetricConf {
		var endpointURL *url.URL
		var err error

		if absoluteUrlRegex.MatchString(metric.Endpoint) {
			endpointURL, err = url.Parse(metric.Endpoint)
		} else {
			endpointURL, err = url.Parse(h.env.ReadUrl + metric.Endpoint)
		}

		if err != nil {
			errorLogger.Printf("Cannot parse url [%v], error: [%v]", metric.Endpoint, err.Error())
			continue
		}

		wg.Add(1)
		go checkServiceReachable(metric.Alias, endpointURL.String(), h.client, hcErrs, &wg)
	}

	wg.Wait()
	close(hcErrs)
	for err := range hcErrs {
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func buildFtHealthcheckUrl(serviceUrl string) (string, error) {
	parsedURL, err := url.Parse(serviceUrl)
	if err != nil {
		return "", err
	}

	var newPath string
	if strings.HasPrefix(parsedURL.Path, "/__") {
		newPath = strings.SplitN(parsedURL.Path[1:], "/", 2)[0] + "/__health"
	} else {
		newPath = "/__health"
	}

	parsedURL.Path = newPath
	return parsedURL.String(), nil
}

func buildAwsHealthcheckUrl(serviceUrl string) (string, error) {
	return serviceUrl + "healthCheckDummyFile", nil
}
