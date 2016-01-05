package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	. "github.com/Financial-Times/publish-availability-monitor/content"
	"github.com/gorilla/mux"
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
	isMarkedDeleted bool
}

// MetricConfig is the configuration of a PublishMetric
type MetricConfig struct {
	Granularity  int      `json:"granularity"` //how we split up the threshold, ex. 120/12
	Endpoint     string   `json:"endpoint"`
	ContentTypes []string `json:"contentTypes"` //list of valid eom types for this metric
	Alias        string   `json:"alias"`
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

	go enableHealthchecks()
	startAggregator()
	readMessages()
}

func enableHealthchecks() {

	healthcheck := &Healthcheck{http.Client{}, *appConfig}
	router := mux.NewRouter()
	router.HandleFunc("/__health", healthcheck.checkHealth())
	router.HandleFunc("/__gtg", healthcheck.gtg)
	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panicf("Couldn't set up HTTP listener: %+v\n", err)
	}
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

	content, err := UnmarshalContent(msg)
	if err != nil {
		warn.Printf("Cannot unmarshal message [%v], error: [%v]", tid, err.Error())
		return err
	}

	uuid := content.GetUUID()
	if !content.IsValid() {
		info.Printf("Message [%v] with UUID [%v] is INVALID, skipping...", tid, uuid)
		return nil
	}

	info.Printf("Message [%v] with UUID [%v] is VALID.", tid, uuid)

	publishDateString := msg.Headers["Message-Timestamp"]
	publishDate, err := time.Parse(dateLayout, publishDateString)
	if err != nil {
		log.Printf("Cannot parse publish date [%v] from message [%v], error: [%v]",
			publishDateString, tid, err.Error())
		return nil
	}

	if isMessagePastPublishSLA(publishDate, appConfig.Threshold) {
		info.Printf("Message [%v] with UUID [%v] is past publish SLA, skipping.", tid, uuid)
		return nil
	}

	scheduleChecks(content, publishDate, tid, content.IsMarkedDeleted())

	// for images we need to check their corresponding image sets
	// the image sets don't have messages of their own so we need to create one
	if content.GetType() == "Image" {
		eomFile, ok := content.(EomFile)
		if !ok {
			log.Printf("Cannot assert that message [%v] with UUID [%v] and type 'Image' is an EomFile.", tid, uuid)
			return nil
		}
		imageSetEomFile := spawnImageSet(eomFile)
		if imageSetEomFile.UUID != "" {
			scheduleChecks(imageSetEomFile, publishDate, tid, false)
		}
	}

	return nil
}

func spawnImageSet(imageEomFile EomFile) EomFile {
	imageSetEomFile := imageEomFile
	imageSetEomFile.Type = "ImageSet"

	imageUUID, err := NewUUIDFromString(imageEomFile.UUID)
	if err != nil {
		warn.Printf("Cannot generate UUID from image UUID string [%v]: [%v], skipping image set check.",
			imageEomFile.UUID, err.Error())
		return EomFile{}
	}

	imageSetUUID, err := GenerateImageSetUUID(*imageUUID)
	if err != nil {
		warn.Printf("Cannot generate image set UUID: [%v], skipping image set check",
			err.Error())
		return EomFile{}
	}

	imageSetEomFile.UUID = imageSetUUID.String()
	return imageSetEomFile
}

func isMessagePastPublishSLA(date time.Time, threshold int) bool {
	passedSLA := date.Add(time.Duration(threshold) * time.Second)
	return time.Now().After(passedSLA)
}

func isSyntheticMessage(tid string) bool {
	return strings.HasPrefix(tid, "SYNTHETIC")
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
