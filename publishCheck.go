package main

import "net/http"

type PublishCheck struct {
	Metric        PublishMetric
	Threshold     int
	CheckInterval int
	ResultSink    chan PublishMetric
}

func NewPublishCheck(pm PublishMetric, t int, ci int, rs chan PublishMetric) *PublishCheck {
	return &PublishCheck{pm, t, ci, rs}
}

func (pc PublishCheck) DoCheck() bool {
	info.Printf("Running check for UUID [%v]\n", pc.Metric.UUID)
	resp, err := http.Head(pc.Metric.endpoint.String() + pc.Metric.UUID)
	if err == nil {
		return resp.StatusCode == 200
	}
	return false
}
