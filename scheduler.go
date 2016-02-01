package main

import (
	"net/url"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/content"
	"fmt"
)

func scheduleChecks(contentToCheck content.Content, publishDate time.Time, tid string, isMarkedDeleted bool) {
	for _, metric := range appConfig.MetricConf {
		endpointURL, err := url.Parse(metric.Endpoint)
		if err != nil {
			errorLogger.Printf("Cannot parse url [%v], error: [%v]", metric.Endpoint, err.Error())
			continue
		}
		if !validType(metric.ContentTypes, contentToCheck.GetType()) {
			continue
		}

		var publishMetric = PublishMetric{
			contentToCheck.GetUUID(),
			false,
			publishDate,
			appConfig.Platform,
			Interval{},
			metric,
			*endpointURL,
			tid,
			isMarkedDeleted,
		}

		var checkInterval = appConfig.Threshold / metric.Granularity
		var publishCheck = NewPublishCheck(publishMetric, appConfig.Threshold, checkInterval, metricSink)
		go scheduleCheck(*publishCheck)
	}
}

func scheduleCheck(check PublishCheck) {

	//the date the SLA expires for this publish event
	publishSLA := check.Metric.publishDate.Add(time.Duration(check.Threshold) * time.Second)

	//compute the actual seconds left until the SLA to compensate for the
	//time passed between publish and the message reaching this point
	secondsUntilSLA := publishSLA.Sub(time.Now()).Seconds()
	infoLogger.Printf("Checking %s. [%v] seconds until SLA.", loggingContextForCheck(check.Metric.config.Alias, check.Metric.UUID, check.Metric.tid), int(secondsUntilSLA))

	//used to signal the ticker to stop after the threshold duration is reached
	quitChan := make(chan bool)
	go func() {
		<-time.After(time.Duration(secondsUntilSLA) * time.Second)
		close(quitChan)
	}()

	secondsSincePublish := time.Since(check.Metric.publishDate).Seconds()
	infoLogger.Printf("Checking %s. [%v] seconds elapsed since publish.", loggingContextForCheck(check.Metric.config.Alias, check.Metric.UUID, check.Metric.tid), int(secondsSincePublish))

	elapsedIntervals := secondsSincePublish / float64(check.CheckInterval)
	infoLogger.Printf("Checking %s. Skipping first [%v] checks", loggingContextForCheck(check.Metric.config.Alias, check.Metric.UUID, check.Metric.tid), int(elapsedIntervals))

	checkNr := int(elapsedIntervals) + 1
	// ticker to fire once per interval
	tickerChan := time.NewTicker(time.Duration(check.CheckInterval) * time.Second)
	for {
		checkSuccessful, ignoreCheck := check.DoCheck()
		if ignoreCheck {
			infoLogger.Printf("Ignore check for %s", loggingContextForCheck(check.Metric.config.Alias, check.Metric.UUID, check.Metric.tid))
			tickerChan.Stop()
			return
		}
		if checkSuccessful {
			tickerChan.Stop()
			check.Metric.publishOK = true

			lower := (checkNr - 1) * check.CheckInterval
			upper := checkNr * check.CheckInterval
			check.Metric.publishInterval = Interval{lower, upper}

			check.ResultSink <- check.Metric
			return
		}
		checkNr++
		select {
		case <-tickerChan.C:
			continue
		case <-quitChan:
			tickerChan.Stop()
			//if we get here, checks were unsuccessful
			check.Metric.publishOK = false
			check.ResultSink <- check.Metric
			return
		}
	}

}
func validType(validTypes []string, eomType string) bool {
	for _, t := range validTypes {
		if t == eomType {
			return true
		}
	}
	return false
}


func loggingContextForCheck(checkType string, uuid string, transactionId string) (string) {
   return fmt.Sprintf("checkType=[%v], uuid=[%v], transaction_id=[%v]", checkType, uuid, transactionId)
}