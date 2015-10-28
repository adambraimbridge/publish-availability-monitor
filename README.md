# publish-availability-monitor
Monitors publish availability and collects related metrics

# Usage
`go get github.com/Financial-Times/publish-availability-monitor`

`publish-availability-monitor -config=int-config.json`

With Docker:

`docker build -t coco/publish-availability-monitor .`

`docker run -ti --env QUEUE_ADDR=<addr> --env URLS=<endpoint1>,<endpoint2> coco/publish-availability-monitor`

# TODO
* logging
* document Graphite limitations (1 datapoint / second)
