# publish-availability-monitor
Monitors publish availability and collects related metrics

# Usage
`go get github.com/Financial-Times/publish-availability-Monitors`

`publish-availability-monitor -config=int-config.json`

With Docker:

`docker build -t coco/publish-availability-monitor .`

`docker run -ti --env QUEUE_ADDR=<addr> --env URL1=<endpoint1>
--env URL2=<endpoint2> --env URL3=<endpoint3> --env URL4=<endpoint4> coco/publish-availability-monitor`
