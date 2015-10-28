package main

type MetricDestination interface {
	Send(pm PublishMetric)
}

type Aggregator struct {
	publishMetricSource       chan PublishMetric
	publishMetricDestinations []MetricDestination
}

func NewAggregator(inputChannel chan PublishMetric, destinations []MetricDestination) *Aggregator {
	return &Aggregator{inputChannel, destinations}
}

func (a *Aggregator) Run() {
	for publishMetric := range a.publishMetricSource {
		for _, sender := range a.publishMetricDestinations {
			go sender.Send(publishMetric)
		}
	}
}
