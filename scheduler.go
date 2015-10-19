package scheduler

type Check struct {
	endpoint //etc.
}

func Schedule() {
	//foreach endpoint, create a new check which will poll that endpoint
	//for the given uuid, every x seconds (granularity) until we reach the threshold (SLA)
	//if the endpoint returns 200 or the threshold is exceeded stop, add result to result channel
	//the aggregator should be on the other side of the channel
}

func NewScheduler() {

}
