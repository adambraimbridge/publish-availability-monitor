package main

// MetricDestination is the interface which defines a method to send
// PublishMetrics to a certain destination.
type MetricDestination interface {
	Send(pm PublishMetric)
}

// Aggregator reads PublishMetrics from a channel and distributes them to
// the configured MetricDestinations.
type Aggregator struct {
	publishMetricSource       chan PublishMetric
	publishMetricDestinations []MetricDestination
}

// NewAggregator returns an Aggregator which reads messages from inputChannel and
// distributes them to destinations.
func NewAggregator(inputChannel chan PublishMetric, destinations []MetricDestination) *Aggregator {
	return &Aggregator{inputChannel, destinations}
}

// Run reads PublishMetrics from a channel and distributes them to a list of
// MetricDestinations.
// Stops reading when the channel is closed.
func (a *Aggregator) Run() {
	for publishMetric := range a.publishMetricSource {
		for _, sender := range a.publishMetricDestinations {
			go sender.Send(publishMetric)
		}
	}
}
