package main

import "github.com/kr/pretty"

type SplunkFeeder struct{}

func (gf SplunkFeeder) Send(pm PublishMetric) {
	info.Printf("Sending the metric to splunk: [%v]", pretty.Formatter(pm))
}
