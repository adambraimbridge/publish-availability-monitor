package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hoisie/mustache"
)

type event struct {
	PublishOK bool //did it meet the SLA?
	Age       time.Duration
	Tid       string
	UUID      string
	Endpoint  string
}

var template = `<html>
	<head>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" />
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"/>
		<title>Publish History</title>
	</head>
	<body>
		<table class="table table-striped">
			<thead>
				<tr>
					<th>Status</th>
					<th>Time</th>
					<th>Transaction ID</th>
					<th>UUID</th>
					<th>Check Type</th>
			</thead>
			<tbody>
		    	{{#events}}
				<tr>
					<td><span class="fa fa-{{#PublishOK}}check{{/PublishOK}}{{^PublishOK}}times{{/PublishOK}}"> </span></td>
					<td>{{Age}} ago</td>
					<td>{{Tid}}</td>
					<td>{{UUID}}</td>
					<td>{{Endpoint}}</td>
				</tr>
				{{/events}}
			</tbody>
		</table>
	</body>
</html>`

func loadHistory(w http.ResponseWriter, r *http.Request) {
	metricContainer.RLock()
	defer metricContainer.RUnlock()

	if strings.Contains(r.Header.Get("Accept"), "text/plain") {
		writePlainTextHistory(w)
	} else {
		writeMustacheHistory(w)
	}
}

func writePlainTextHistory(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/plain")
	for i := len(metricContainer.publishMetrics) - 1; i >= 0; i-- {
		fmt.Fprintf(w, "%d. %v\n\n", len(metricContainer.publishMetrics)-i, metricContainer.publishMetrics[i])
	}
}

func writeMustacheHistory(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/html")

	sz := len(metricContainer.publishMetrics)
	events := make([]event, 0, sz)
	for i := sz - 1; i >= 0; i-- {
		events = append(events,
			event{
				metricContainer.publishMetrics[i].publishOK,
				time.Since(metricContainer.publishMetrics[i].publishDate),
				metricContainer.publishMetrics[i].tid,
				metricContainer.publishMetrics[i].UUID,
				metricContainer.publishMetrics[i].config.Alias,
			})
	}

	ctx := make(map[string]interface{})
	ctx["events"] = events

	html := mustache.Render(template, ctx)
	fmt.Fprintln(w, html)
}
