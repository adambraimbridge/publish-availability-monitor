package main

import (
	"net/url"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/content"
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
	infoLogger.Printf("Seconds until SLA for [%v] : [%v]", check.Metric.UUID, int(secondsUntilSLA))

	//used to signal the ticker to stop after the threshold duration is reached
	quitChan := make(chan bool)
	go func() {
		<-time.After(time.Duration(secondsUntilSLA) * time.Second)
		close(quitChan)
	}()

	secondsSincePublish := time.Since(check.Metric.publishDate).Seconds()
	infoLogger.Printf("Seconds elapsed since publish for [%v] : [%v]", check.Metric.UUID, int(secondsSincePublish))

	elapsedIntervals := secondsSincePublish / float64(check.CheckInterval)
	infoLogger.Printf("Skipping first [%v] checks for [%v]", int(elapsedIntervals), check.Metric.UUID)

	checkNr := int(elapsedIntervals) + 1
	// ticker to fire once per interval
	tickerChan := time.NewTicker(time.Duration(check.CheckInterval) * time.Second)
	for {
		if check.DoCheck() {
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
