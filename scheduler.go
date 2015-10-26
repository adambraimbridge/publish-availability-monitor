package main

import (
	"log"
	"net/url"
	"time"
)

func scheduleChecks(eomFile EomFile, publishDate time.Time) {
	for _, conf := range appConfig.MetricConf {
		endpointUrl, err := url.Parse(conf.Endpoint)
		if err != nil {
			log.Printf("Cannot parse url [%v], error: [%v]", conf.Endpoint, err.Error())
			continue
		}
		if conf.ContentType != "" && conf.ContentType != eomFile.Type {
			continue
		}
		var publishMetric = PublishMetric{
			eomFile.UUID,
			false,
			publishDate,
			appConfig.Platform,
			Interval{},
			conf,
			*endpointUrl,
		}

		var checkInterval = appConfig.Threshold / conf.Granularity
		var publishCheck = NewPublishCheck(publishMetric, appConfig.Threshold, checkInterval, metricSink)
		go scheduleCheck(*publishCheck)
	}
}

func scheduleCheck(check PublishCheck) {

	quitChan := make(chan bool)
	checkNr := 1

	//used to signal the ticker to stop after the threshold duration is reached
	go func() {
		<-time.After(time.Duration(check.Threshold) * time.Second)
		close(quitChan)
	}()

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
			return
		}
	}

	//if we get here, checks were unsuccessful
	check.Metric.publishOK = false
	check.ResultSink <- check.Metric
}
