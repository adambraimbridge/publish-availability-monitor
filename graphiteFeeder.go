package main

import (
	"fmt"
	"net"
	"strconv"
)

type GraphiteFeeder struct {
	host       string
	port       int
	connection net.Conn
}

const prefix = "high-res.content.publish.metrics."
const durationSuffix = "duration %d %d\n"
const suffix = " %d %d\n"

func NewGraphiteFeeder(host string, port int) *GraphiteFeeder {
	connection, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))
	//TODO handle error
	return &GraphiteFeeder{host, port, connection}
}

func (gf GraphiteFeeder) Send(pm PublishMetric) {
	metricCommon := prefix + appConfig.Platform + "." + pm.config.Alias + "."
	duration := metricCommon + "duration" + suffix
	status := metricCommon + "status" + suffix

	_, err := fmt.Fprintf(gf.connection, duration, pm.publishInterval.upperBound, pm.publishDate.Unix())
	if err != nil {
		warn.Printf("Error sending stuff to graphite: [%v]", err.Error())
		//TODO re-establish connection, retry
	}

	_, err = fmt.Fprintf(gf.connection, status, Btoi(pm.publishOK), pm.publishDate.Unix())
	if err != nil {
		warn.Printf("Error sending stuff to graphite: [%v]", err.Error())
		//TODO re-establish connection, retry
	}
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
