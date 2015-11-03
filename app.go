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
	ContentType string `json:"contentType"`
	Alias       string `json:"alias"`
}

type GraphiteConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SplunkConfig struct {
	FilePath string `json:"logFilePath"`
}

type AppConfig struct {
	Threshold    int                  `json:"threshold"` //pub SLA in seconds, ex. 120
	QueueConf    consumer.QueueConfig `json:"queueConfig"`
	MetricConf   []MetricConfig       `json:"metricConfig"`
	Platform     string               `json:"platform"`
	GraphiteConf GraphiteConfig       `json:"graphite-config"`
	SplunkConf   SplunkConfig         `json:"splunk-config"`
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

func main() {
	initLogs(os.Stdout, os.Stdout, os.Stderr)
	flag.Parse()

	var err error
	appConfig, err = ParseConfig(*configFileName)
	if err != nil {
		log.Printf("Cannot load configuration: [%v]", err)
		return
	}

	startAggregator()
	readMessages()
}

func readMessages() {
	iterator := consumer.NewIterator(appConfig.QueueConf)
	for {
		msgs, err := iterator.NextMessages()
		if err != nil {
			warn.Printf("Could not read messages: [%v]", err.Error())
			continue
		}
		for _, m := range msgs {
			go PublishMessageListener{}.OnMessage(m)
		}
	}
}

func startAggregator() {
	var destinations []MetricDestination

	if len(appConfig.GraphiteConf.Host) != 0 && appConfig.GraphiteConf.Port != 0 {
		graphiteFeeder := NewGraphiteFeeder(appConfig.GraphiteConf.Host, appConfig.GraphiteConf.Port)
		destinations = append(destinations, graphiteFeeder)
	}
	splunkFeeder := NewSplunkFeeder(appConfig.SplunkConf.FilePath)
	destinations = append(destinations, splunkFeeder)
	aggregator := NewAggregator(metricSink, destinations)
	go aggregator.Run()
}

func (listener PublishMessageListener) OnMessage(msg consumer.Message) error {
	tid := msg.Headers["X-Request-Id"]
	info.Printf("Received message with TID [%v]", tid)

	if isSyntheticMessage(tid) {
		info.Printf("Message [%v] is INVALID: synthetic, skipping...", tid)
		return nil
	}

	if !isMessageValid(msg) {
		info.Printf("Message [%v] is INVALID, skipping...", tid)
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

	info.Printf("Message [%v] is VALID.", tid)

	publishDateString := msg.Headers["Message-Timestamp"]
	publishDate, err := time.Parse(dateLayout, publishDateString)
	if err != nil {
		log.Printf("Cannot parse publish date [%v] from message [%v], error: [%v]",
			publishDateString, tid, err.Error())
		return nil
	}

	if isMessagePastPublishSLA(publishDate, appConfig.Threshold) {
		info.Printf("Message [%v] is past publish SLA, skipping.", tid)
		return nil
	}

	scheduleChecks(eomFile, publishDate, tid)
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
