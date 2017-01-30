package history

import (
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/Financial-Times/publish-availability-monitor/checks"
	"github.com/aymerick/raymond"
)

type PublishStatus int

const (
	IGNORED PublishStatus = iota
	INVALID
	CHECKING
	FAILED
	SUCCESS
)

type publishCheck struct {
	UUID        string
	Type        string
	Environment string
	Status      PublishStatus
}

type publishEvent struct {
	At     string
	Tid    string
	Status PublishStatus
	Checks []*publishCheck
}

var (
	lock    sync.Mutex
	history []*publishEvent

	template string
)

func init() {
	history = make([]*publishEvent, 0)

	by, err := ioutil.ReadFile("./resources/history.hbs")
	if err != nil {
		infoLogger.Fatalf("unable to read template: %v", err)
	}

	template = string(by)
}

func MonitorPublish(tid string) {
	lock.Lock()
	defer lock.Unlock()

	p := publishEvent{time.Now().Format(time.RFC3339), tid, CHECKING, make([]*publishCheck, 0)}
	history = append(history, &p)
}

func MonitorPublishCheck(tid string, uuid string, checkType string, environment string) {
	c := publishCheck{uuid, checkType, environment, CHECKING}

	lock.Lock()
	defer lock.Unlock()

	for _, p := range history {
		if p.Tid == tid {
			p.Checks = append(p.Checks, &c)
			return
		}
	}

	infoLogger.Printf("No publishEvent found for %s", tid)
}

func HandlePublishResult(tid string, status PublishStatus) {
	lock.Lock()
	defer lock.Unlock()

	for _, p := range history {
		if p.Tid == tid {
			p.Status = status
			break
		}
	}
}

func HandlePublishCheckResult(tid string, uuid string, checkType string, environment string, status PublishStatus) {
	lock.Lock()
	defer lock.Unlock()

	for _, p := range history {
		if p.Tid == tid {
			allOK := true
			allIgnored := true
			for _, c := range p.Checks {
				if c.UUID == uuid && c.Type == checkType && c.Environment == environment {
					c.Status = status
				}

				// publish status is FAILED if any check status is FAILED
				if c.Status == FAILED {
					p.Status = FAILED
				}
				// publish status is SUCCESS/IGNORED if all check statuses are SUCCESS/IGNORED
				allOK = allOK && (c.Status == SUCCESS)
				allIgnored = allIgnored && (c.Status == IGNORED)
			}

			if allOK {
				p.Status = SUCCESS
			}
			if allIgnored {
				p.Status = IGNORED
			}

			break
		}
	}
}

func WriteHistory(out io.Writer) {
	sz := len(history)
	events := make([]publishEvent, 0, sz)
	for i := sz - 1; i >= 0; i-- {
		if !isIgnorable(*history[i]) {
			events = append(events, *history[i])
		}
	}

	ctx := make(map[string]interface{})
	ctx["events"] = events

	html, err := raymond.Render(template, ctx)
	if err == nil {
		//	html := mustache.Render(template, ctx)
		fmt.Fprintln(out, html)
	} else {
		fmt.Fprintln(out, err.Error())
	}

}

func Forget(tid string) {
	infoLogger.Printf("forgetting %s", tid)
	lock.Lock()
	defer lock.Unlock()

	for i, p := range history {
		if p.Tid == tid {
			tail := history[i+1:]
			history = append(history[:i], tail...)
			break
		}
	}
}

func isIgnorable(event publishEvent) bool {
	return (event.Status == IGNORED) || checks.IsIgnorableMessage(event.Tid)
}
