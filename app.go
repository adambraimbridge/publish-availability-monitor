package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"fmt"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/Financial-Times/publish-availability-monitor/content"
	"github.com/Financial-Times/publish-availability-monitor/feeds"
	"github.com/Financial-Times/publish-availability-monitor/logformat"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	log "github.com/Sirupsen/logrus"
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
	Health       string   `json:"health,omitempty"`
	ApiKey       string   `json:"apiKey,omitempty"`
}

// SplunkConfig holds the SplunkFeeder-specific configuration
type SplunkConfig struct {
	LogPrefix string `json:"logPrefix"`
}

// AppConfig holds the application's configuration
type AppConfig struct {
	Threshold           int                  `json:"threshold"` //pub SLA in seconds, ex. 120
	QueueConf           consumer.QueueConfig `json:"queueConfig"`
	MetricConf          []MetricConfig       `json:"metricConfig"`
	SplunkConf          SplunkConfig         `json:"splunk-config"`
	HealthConf          HealthConfig         `json:"healthConfig"`
	ValidationEndpoints map[string]string    `json:"validationEndpoints"` //contentType to validation endpoint mapping, ex. { "EOM::Story": "http://methode-article-transformer/content-transform" }
}

// HealthConfig holds the application's healthchecks configuration
type HealthConfig struct {
	FailureThreshold int `json:"failureThreshold"`
}

// Environment defines an environment in which the publish metrics should be checked
type Environment struct {
	Name     string `json:"name"`
	ReadUrl  string `json:"read-url"`
	S3Url    string `json:"s3-url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credentials struct {
	EnvName  string `json:"env-name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type publishHistory struct {
	sync.RWMutex
	publishMetrics []PublishMetric
}

const dateLayout = time.RFC3339Nano

var configFileName = flag.String("config", "", "Path to configuration file")

var etcdPeers = flag.String("etcd-peers", "http://localhost:2379", "Comma-separated list of addresses of etcd endpoints to connect to")
var etcdReadEnvKey = flag.String("etcd-read-env-key", "/ft/config/monitoring/read-urls", "etcd key that lists the read environment URLs")
var etcdS3EnvKey = flag.String("etcd-s3-env-key", "/ft/config/monitoring/s3-image-bucket-urls", "etcd key that lists the S3 image bucket URLs")
var etcdCredKey = flag.String("etcd-cred-key", "/ft/_credentials/publish-read/read-credentials", "etcd key that lists the read environment credentials")
var etcdValidatorCredKey = flag.String("etcd-validator-cred-key", "/ft/_credentials/publish-read/validator-credentials", "etcd key that specifies the validator credentials")

var envsFileName = flag.String("envs-file-name", "/etc/pam/envs/read-environments.json", "Path to json file that contains environments configuration")
var envCredentialsFileName = flag.String("envs-credentials-file-name", "/etc/pam/credentials/read-environments-credentials.json", "Path to json file that contains environments credentials")
var validatorCredentialsFileName = flag.String("validator-credentials-file-name", "/etc/pam/credentials/validator-credentials.json", "Path to json file that contains validation endpoints configuration")
var configRefreshPeriod = flag.Int("config-refresh-period", 1, "Refresh period for configuration in minutes. By default it is 1 minute.")

var appConfig *AppConfig
var environments = make(map[string]Environment)
var subscribedFeeds = make(map[string][]feeds.Feed)
var metricSink = make(chan PublishMetric)
var metricContainer publishHistory
var validatorCredentials string
var configFilesHashValues = make(map[string]string)
var carouselTransactionIDRegExp = regexp.MustCompile(`^.+_carousel_[\d]{10}.*$`)

func init() {
	f := logformat.NewSLF4JFormatter(`.*/github\.com/Financial-Times/.*`)
	log.SetFormatter(f)
}

func main() {
	flag.Parse()

	var err error
	appConfig, err = ParseConfig(*configFileName)
	if err != nil {
		log.Fatal("Cannot load configuration: [%v]", err)
		return
	}


	log.Info("EtcdPeers value: [%s]", *etcdPeers)
	if len(*etcdPeers) > 0{
		log.Info("Sourcing dynamic configs from ETCD")
		go DiscoverEnvironmentsAndValidators(etcdPeers, etcdReadEnvKey, etcdCredKey, etcdS3EnvKey, etcdValidatorCredKey, environments)
	} else{
		log.Info("Sourcing dynamic configs from file")
		err = updateEnvsIfChanged(*envsFileName, *envCredentialsFileName)
		if err != nil {
			panic(fmt.Sprintf("Cannot load envs config, error was: [%v]", err))
		}

		err = updateValidationCredentialsIfChanged(*validatorCredentialsFileName)
		if err != nil {
			panic(fmt.Sprintf("Cannot load validation credentials, error was: [%v]", err))
		}

		go watchConfigFiles(*envsFileName, *envCredentialsFileName, *validatorCredentialsFileName, *configRefreshPeriod)
	}

	metricContainer = publishHistory{sync.RWMutex{}, make([]PublishMetric, 0)}

	go startHttpListener()

	startAggregator()
	readMessages()
}

func startHttpListener() {
	router := mux.NewRouter()
	setupHealthchecks(router)
	router.HandleFunc("/__history", loadHistory)

	router.HandleFunc(status.PingPath, status.PingHandler)
	router.HandleFunc(status.PingPathDW, status.PingHandler)

	router.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)
	router.HandleFunc(status.BuildInfoPathDW, status.BuildInfoHandler)

	attachProfiler(router)

	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panicf("Couldn't set up HTTP listener: %+v\n", err)
	}
}

func setupHealthchecks(router *mux.Router) {
	hc := newHealthcheck(appConfig, &metricContainer)
	router.HandleFunc("/__health", hc.checkHealth)
	router.HandleFunc("/__gtg", status.NewGoodToGoHandler(hc.GTG))
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

func readMessages() {
	c := consumer.NewConsumer(appConfig.QueueConf, handleMessage, &http.Client{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		c.Start()
		wg.Done()
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	c.Stop()
	wg.Wait()
}

func startAggregator() {
	var destinations []MetricDestination

	splunkFeeder := NewSplunkFeeder(appConfig.SplunkConf.LogPrefix)
	destinations = append(destinations, splunkFeeder)
	aggregator := NewAggregator(metricSink, destinations)
	go aggregator.Run()
}

func loadHistory(w http.ResponseWriter, r *http.Request) {
	metricContainer.RLock()
	for i := len(metricContainer.publishMetrics) - 1; i >= 0; i-- {
		fmt.Fprintf(w, "%d. %v\n\n", len(metricContainer.publishMetrics)-i, metricContainer.publishMetrics[i])
	}
	metricContainer.RUnlock()
}

func handleMessage(msg consumer.Message) {
	tid := msg.Headers["X-Request-Id"]
	log.Infof("Received message with TID [%v]", tid)

	if isIgnorableMessage(tid) {
		log.Infof("Message [%v] is ignorable. Skipping...", tid)
		return
	}

	publishDateString := msg.Headers["Message-Timestamp"]
	publishDate, err := time.Parse(dateLayout, publishDateString)
	if err != nil {
		log.Errorf("Cannot parse publish date [%v] from message [%v], error: [%v]",
			publishDateString, tid, err.Error())
		return
	}

	publishedContent, err := content.UnmarshalContent(msg)
	if err != nil {
		log.Warnf("Cannot unmarshal message [%v], error: [%v]", tid, err.Error())
		return
	}

	var paramsToSchedule []*schedulerParam

	for _, preCheck := range mainPreChecks() {
		ok, scheduleParam := preCheck(publishedContent, tid, publishDate)
		if ok {
			paramsToSchedule = append(paramsToSchedule, scheduleParam)
		} else {
			//if a main check is not ok, additional checks make no sense
			return
		}
	}

	for _, preCheck := range additionalPreChecks() {
		ok, scheduleParam := preCheck(publishedContent, tid, publishDate)
		if ok {
			paramsToSchedule = append(paramsToSchedule, scheduleParam)
		}
	}

	for _, scheduleParam := range paramsToSchedule {
		scheduleChecks(scheduleParam)
	}
}

func isIgnorableMessage(tid string) bool {
	return isSyntheticTransactionID(tid) || isContentCarouselTransactionID(tid)
}

func isSyntheticTransactionID(tid string) bool {
	return strings.HasPrefix(tid, "SYNTHETIC")
}

func isContentCarouselTransactionID(tid string) bool {
	return carouselTransactionIDRegExp.MatchString(tid)
}

func (pm PublishMetric) String() string {
	return fmt.Sprintf("Tid: %s, UUID: %s, Platform: %s, Endpoint: %s, PublishDate: %s, Duration: %d, Succeeded: %t.",
		pm.tid,
		pm.UUID,
		pm.platform,
		pm.config.Alias,
		pm.publishDate.String(),
		pm.publishInterval.upperBound,
		pm.publishOK,
	)

}
