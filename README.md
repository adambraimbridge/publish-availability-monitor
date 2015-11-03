# publish-availability-monitor

[![Circle CI](https://circleci.com/gh/Financial-Times/publish-availability-monitor/tree/master.svg?style=svg)](https://circleci.com/gh/Financial-Times/publish-availability-monitor/tree/master)

Monitors publish availability and collects related metrics

# Usage
`go get github.com/Financial-Times/publish-availability-monitor`

`publish-availability-monitor -config=int-config.json`

With Docker:

`docker build -t coco/publish-availability-monitor .`

`docker run -ti --env QUEUE_ADDR=<addr> --env URLS=<endpoint1>,<endpoint2> coco/publish-availability-monitor`

# TODO
* logging

# Puppet module
If you want a new puppet module version, you have to manually increment the version in `puppet/ft-publish_availability_monitor/Modulefile` AND create a `git tag <version>` corresponding to the module version. This way Jenkins can pick up the tagged version to build & deploy to forge.
