package main

type Aggregator struct {
	publishMetricSource chan PublishMetric
}

func NewAggregator(inputChannel chan PublishMetric) *Aggregator {
	return &Aggregator{inputChannel}
}

func (a *Aggregator) Run() {
	for publishMetric := range a.publishMetricSource {
		info.Printf("Received publish metric: [%v]\n", publishMetric)
	}
}
