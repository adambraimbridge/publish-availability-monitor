# publish-availability-monitor
[![Circle CI](https://circleci.com/gh/Financial-Times/publish-availability-monitor/tree/master.png?style=shield)](https://circleci.com/gh/Financial-Times/publish-availability-monitor/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Financial-Times/publish-availability-monitor)](https://goreportcard.com/report/github.com/Financial-Times/publish-availability-monitor) [![Coverage Status](https://coveralls.io/repos/github/Financial-Times/publish-availability-monitor/badge.svg?branch=master)](https://coveralls.io/github/Financial-Times/publish-availability-monitor?branch=master) [![codecov](https://codecov.io/gh/Financial-Times/publish-availability-monitor/branch/master/graph/badge.svg)](https://codecov.io/gh/Financial-Times/publish-availability-monitor)

Monitors publish availability and collects related metrics. Collected metrics are sent to various systems (ex. Splunk).

# Usage
```
go get -u github.com/kardianos/govendor
go get github.com/Financial-Times/publish-availability-monitor
govendor sync
publish-availability-monitor -config=int-config.json
```

With Docker:

`docker build -t coco/publish-availability-monitor .`

`docker run -ti --env QUEUE_ADDR=<addr> --env S3_URL=<S3 bucket URL> --env CONTENT_URL=<document store api article endpoint path> --env LISTS_URL=<document store api lists endpoint path> --env NOTIFICATIONS_URL=<notifications read path> --env NOTIFICATIONS_PUSH_URL=<notifications push path> --env METHODE_ARTICLE_TRANSFORMER_URL=<methode article transformer URL>  --env METHODE_CONTENT_PLACEHOLDER_MAPPER_URL=<methode content placeholder mapper URL> coco/publish-availability-monitor`

# Build and deploy
__Note that deployment to FTP2 is no longer supported.__
* Tagging a release in GitHub triggers a build in DockerHub.

# Endpoint Check Configuration

```
//this is the SLA for content publish, in seconds
//the app will check for content availability until this threshold is reached
//or the content was found  
"threshold": 120,
```

```
//Configuration for the queue we will read from
"queueConfig": {
	//URL(s) of the messaging system
	"address": [queue address1, queue address 2],
	//The group name by which this app will be knows to the messaging system
	"group": "YourGroupName",
	//The topic we want to get messages from
	"topic": "YourTopic",
	//the name of the queue we want to use
	"queue": "yourQueue"
},
```

```
//for each endpoint we want to check content against, we need a metricConfig entry
"metricConfig": [
{
	//the path (or absolute URL) of the endpoint.
	//paths are resolved relative to the read_url setting (see below) on a per-environment basis
	//to check content, the UUID is appended at the end of this URL
	//should end with /
	"endpoint": "endpointURL",
	//used to associate endpoint-specific behavior with the endpoint
	//each alias should have an entry in the endpointSpecificChecks map
	"alias": "content",
	//defines how often we check this endpoint
	//the check interval is threshold / granularity
	//in this case, 120 / 40 = 3 -> we check every 3 seconds
	"granularity": 40
},
{
	"endpoint": "endpointURL",
	"granularity": 40,
	"alias": "S3",
	//optional field to indicate that this endpoint should only be checked
	//for content of a certain type
	//if not present, all content will be checked against this endpoint
	"contentTypes": ["Image"]
}
],
```

```
//feeder-specific configuration
//for each feeder, we need a new struct, new field in AppConfig for it, and
//handling for the feeder in startAggregator()
"splunk-config": {
	"logFilePath": "/var/log/apps/pam.log"
}
```

# Environment Configuration
The app checks environments configuration as well as validation credentials every minute (configurable) and it reloads them if changes are detected.

## JSON example for environments configuration:
 <pre>
     [
       {
         "name":"pre-prod-uk",
         "read-url": "https://pre-prod-uk.ft.com",
         "s3-url": "http://com.ft.imagepublish.amazonaws.com"
       },
       {
         "name":"pre-prod-us",
         "read-url": "https://pre-prod-us.ft.com",
         "s3-url": "http://com.ft.imagepublish.amazonaws.com"
       }       
     ]
 </pre>
## JSON example for environments credentials configuration:
 <pre>
 [
   {
         "env-name": "pre-prod-uk",
         "username": "dummy-username",
         "password": "dummy-pwd"
   },
   {
     "env-name": "pre-prod-us",
     "username": "dummy-username",
     "password": "dummy-pwd"
   }   
       
 ]
  </pre>
## JSON example for validation credentials configuration:
 <pre>
  {
    "username": "dummy-username",
    "password": "dummy-password"
  }
 </pre>
 
Checks that have already been initiated are unaffected by changes to the values above.
