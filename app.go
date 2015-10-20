package main

import (
	"flag"
	"net/url"
	"time"
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

type QueueConfig struct {
	Address string `json:"address"`
	Group   string `json:"group"`
	Topic   string `json:"topic"`
	Queue   string `json:"queue"`
}

type MetricConfig struct {
	Threshold   int `json:"threshold"`   //pub SLA in seconds, ex. 120
	Granularity int `json:"granularity"` //how we split up the threshold, ex. 120/12
}

type AppConfig struct {
	Endpoints  []string     `json:"endpoints"`
	QueueConf  QueueConfig  `json:"queueConfig"`
	MetricConf MetricConfig `json:"metricConfig"`
	//TODO feeder configs
}

func main() {
	//read config (into structs?)
	configFileName := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	ParseConfig(*configFileName)
	/*
		consumer := consumer.NewConsumer(QueueConfig.address etc.)
		scheduler := scheduler.NewScheduler()
		aggregator := aggregator.NewAggregator()
		validator := validator.NewValidator()

		consumer.Consume(/*this?/)
	*/
	//maybe separate the distributor so it just waits for metrics from the aggregator like a servlet?
	//call consume
}

/*
func OnMessage(msg consumer.Message) error {
	//if message is not valid, skip
	//generate timestamp (this is the moment we measure the publish from)
	//scheduler.scheduleChecks(message, publishMetric)
	//connect the scheduler with the aggregator with channels or something
	//so when each scheduler is finished, the aggregator reads the results

}
*/
