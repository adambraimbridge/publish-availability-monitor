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

// Interval is a simple representation of an interval of time, with a lower and
// upper boundary
type Interval struct {
	lowerBound int
	upperBound int
}

// PublishMetric holds the information about the metric we are measuring.
type PublishMetric struct {
	UUID            string
	publishOK       bool      //did it meet the SLA?
	publishDate     time.Time //the time WE get the message
	platform        string
	publishInterval Interval //the interval it was actually published in, ex. (10,20)
	config          MetricConfig
	endpoint        url.URL
	tid             string
}

// MetricConfig is the configuration of a PublishMetric
type MetricConfig struct {
	Granularity int    `json:"granularity"` //how we split up the threshold, ex. 120/12
	Endpoint    string `json:"endpoint"`
	ContentType string `json:"contentType"`
	Alias       string `json:"alias"`
}

// SplunkConfig holds the SplunkFeeder-specific configuration
type SplunkConfig struct {
	LogPrefix string `json:"logPrefix"`
}

// AppConfig holds the application's configuration
type AppConfig struct {
	Threshold  int                  `json:"threshold"` //pub SLA in seconds, ex. 120
	QueueConf  consumer.QueueConfig `json:"queueConfig"`
	MetricConf []MetricConfig       `json:"metricConfig"`
	Platform   string               `json:"platform"`
	SplunkConf SplunkConfig         `json:"splunk-config"`
}

// EomFile models content as it is stored in our CMS
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
			go handleMessage(m)
		}
	}
}

func startAggregator() {
	var destinations []MetricDestination

	splunkFeeder := NewSplunkFeeder(appConfig.SplunkConf.LogPrefix)
	destinations = append(destinations, splunkFeeder)
	aggregator := NewAggregator(metricSink, destinations)
	go aggregator.Run()
}

func handleMessage(msg consumer.Message) error {
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
		info.Printf("Message [%v] with UUID [%v] is INVALID, skipping...", tid, eomFile.UUID)
		return nil
	}

	info.Printf("Message [%v] with UUID [%v] is VALID.", tid, eomFile.UUID)

	publishDateString := msg.Headers["Message-Timestamp"]
	publishDate, err := time.Parse(dateLayout, publishDateString)
	if err != nil {
		log.Printf("Cannot parse publish date [%v] from message [%v], error: [%v]",
			publishDateString, tid, err.Error())
		return nil
	}

	if isMessagePastPublishSLA(publishDate, appConfig.Threshold) {
		info.Printf("Message [%v] with UUID [%v] is past publish SLA, skipping.", tid, eomFile.UUID)
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
