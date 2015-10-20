package main

import (
	"net/url"
	"time"
	consumer "github.com/Financial-Times/go-message-queue-consumer"
)

type Interval struct {
	lowerBound int
	upperBound int
}

type MetricConfig struct {
	threshold   int //pub SLA in seconds, ex. 120
	granularity int //how we split up the threshold, ex. 120/12

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
	address string
	group   string
	topic   string
	queue   string
}

type AppConfig struct {
	endpoints   []url.URL
	queueConfig QueueConfig
	//TODO feeder configs
}

func main() {
	//read config (into structs?)
	consumer := consumer.NewConsumer(QueueConfig.address etc.)
	scheduler := scheduler.NewScheduler()
	aggregator := aggregator.NewAggregator()
	validator := validator.NewValidator()
	
	consumer.Consume(/*this?*/)
	
	//maybe separate the distributor so it just waits for metrics from the aggregator like a servlet?
	//call consume
}

func OnMessage(msg consumer.Message) error {
	//if message is not valid, skip
	//generate timestamp (this is the moment we measure the publish from) 
	//scheduler.scheduleChecks(message, publishMetric) 
	//connect the scheduler with the aggregator with channels or something
	//so when each scheduler is finished, the aggregator reads the results
	
}
