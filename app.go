package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
)

type Interval struct {
	lowerBound int
	upperBound int
}

type PublishMetric struct {
	UUID            string
	publishOK       bool      //did it meet the SLA?
	publishDate     time.Time //the time WE get the message
	platform        string
	publishInterval Interval //the interval it was actually published in, ex. (10,20)
	config          MetricConfig
	endpoint        url.URL
}

type MetricConfig struct {
	Granularity int    `json:"granularity"` //how we split up the threshold, ex. 120/12
	Endpoint    string `json:"endpoint"`
}

type AppConfig struct {
	Threshold  int                  `json:"threshold"` //pub SLA in seconds, ex. 120
	QueueConf  consumer.QueueConfig `json:"queueConfig"`
	MetricConf []MetricConfig       `json:"metricConfig"`
	Platform   string               `json:"platform"`
	//TODO feeder configs
}

type PublishMessageListener struct{}

type EomFile struct {
	UUID             string `json:"uuid"`
	Type             string `json:"type"`
	Value            string `json:"value"`
	Attributes       string `json:"attributes"`
	SystemAttributes string `json:"systemAttributes"`
}

const dateLayout = "2006-01-02T15:04:05.000Z"
const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var info *log.Logger
var warn *log.Logger
var configFileName = flag.String("config", "", "Path to configuration file")
var appConfig *AppConfig
var metricSink = make(chan PublishMetric)
var err error

func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	flag.Parse()

	var err error
	appConfig, err = ParseConfig(*configFileName)
	if err != nil {
		log.Printf("Cannot load configuration: [%v]", err)
		return
	}

	aggregator := NewAggregator(metricSink)
	go aggregator.Run()

	messageConsumer := consumer.NewConsumer(appConfig.QueueConf)
	err = messageConsumer.Consume(PublishMessageListener{}, 8)
	if err != nil {
		log.Printf("Cannot start listening for messages: [%v]", err.Error())
		return
	}

	//maybe separate the distributor so it just waits for metrics from the aggregator like a servlet?
}

func (listener PublishMessageListener) OnMessage(msg consumer.Message) error {
	tid := msg.Headers["X-Request-Id"]
	info.Printf("Received message with TID [%v]", tid)

	if !isMessageValid(msg) {
		return nil
	}

	var eomFile EomFile
	err := json.Unmarshal([]byte(msg.Body), &eomFile)
	if err != nil {
		log.Printf("Cannot unmarshal message [%v], error: [%v]", tid, err.Error())
		return err
	}

	if !isEomfileValid(eomFile) {
		info.Printf("Message [%v] is INVALID, skipping...", tid)
		return nil
	}

	info.Printf("Message [%v] is VALID, scheduling checks...", tid)

	publishDateString := msg.Headers["Message-Timestamp"]
	publishDate, err := time.Parse(dateLayout, publishDateString)
	if err != nil {
		log.Printf("Cannot parse publish date [%v] from message [%v], error: [%v]",
			publishDateString, tid, err.Error())
		return nil
	}

	var publishChecks []PublishCheck
	for _, conf := range appConfig.MetricConf {
		info.Println("MetricConf: %# v", conf)
		endpointUrl, err := url.Parse(conf.Endpoint)
		if err != nil {
			log.Printf("Cannot parse url [%v], error: [%v]", conf.Endpoint, err.Error())
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
		publishChecks = append(publishChecks, *publishCheck)
		go scheduleCheck(*publishCheck)
	}

	//info.Printf("Checks to schedule: %# v", pretty.Formatter(publishChecks))
	return nil
}

func initLogs(infoHandle io.Writer, warnHandle io.Writer, panicHandle io.Writer) {
	//to be used for INFO-level logging: info.Println("foo is now bar")
	info = log.New(infoHandle, "INFO  - ", logPattern)
	//to be used for WARN-level logging: warn.Println("foo is now bar")
	warn = log.New(warnHandle, "WARN  - ", logPattern)

	log.SetFlags(logPattern)
	log.SetPrefix("ERROR - ")
	log.SetOutput(panicHandle)
}

func scheduleCheck(check PublishCheck) {

	quitChan := make(chan bool)

	go func() {
		<-time.After(check.Threshold * time.Second)
		close(quitChan)
	}()

	if check.DoCheck() {
		check.Metric.publishOK = true
		//TODO calculate interval or publishTime
		check.ResultSink <- check.Metric
		return
	}
	// fire once per second
	t := time.NewTicker(check.CheckInterval * time.Second)
	func() {
		for {
			if check.DoCheck() {
				check.Metric.publishOK = true
				//TODO calculate interval or publishTime
				check.ResultSink <- check.Metric
				t.Stop()
				return
			}
			select {
			case <-t.C:
			case <-quitChan:
				t.Stop()
				return
			}
		}
	}()
	check.Metric.publishOK = false
	check.ResultSink <- check.Metric
}
