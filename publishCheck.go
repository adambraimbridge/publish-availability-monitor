package main

import "time"

type PublishCheck struct {
	Metric        PublishMetric
	Threshold     time.Duration
	CheckInterval time.Duration
	ResultSink    chan PublishMetric
}

func NewPublishCheck(pm PublishMetric, t int, ci int, rs chan PublishMetric) *PublishCheck {
	return &PublishCheck{pm, time.Duration(t), time.Duration(ci), rs}
}

func (pc PublishCheck) DoCheck() bool {
	info.Printf("Checking endpoint [%v]", pc.Metric.endpoint)
	return false
}
