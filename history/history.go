package history

import (
	"fmt"
	"io"
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

	template = `<html>
	<head>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" />
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"/>
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
		<style type="text/css">
.republish, .forget {
	cursor: pointer;
}
		</style>
		<meta http-equiv="refresh" content="30"/>
		<title>Publish History</title>
	</head>
	<body>
	<div class="container">
		<table class="table table-striped">
			<thead>
				<tr>
					<th>Time</th>
					<th>Transaction ID</th>
					<th>UUID</th>
					<th>Check Type</th>
					<th>Status</th>
					<th>Actions</th>
			</thead>
		    	{{#events}}
			<tbody>
		    	<tr
					{{#equal Status "1"}}class="warning"{{/equal}}
					{{#equal Status "3"}}class="danger"{{/equal}}
		    	 >
					<td>{{At}}</td>
					<td class="tid" colspan="3">{{Tid}}</td>
					<td><span 
						{{#equal Status "0"}}class="fa fa-ellipsis-h" title="Ignored"{{/equal}}
						{{#equal Status "1"}}class="fa fa-exclamation" title="Invalid"{{/equal}}
						{{#equal Status "2"}}class="fa fa-clock-o" title="In Progress"{{/equal}}
						{{#equal Status "3"}}class="fa fa-times" title="Failed"{{/equal}}
						{{#equal Status "4"}}class="fa fa-check" title="Success"{{/equal}}
						> </span>
					</td>
					<td>
						<button type="button" class="btn btn-link republish"><span class="fa fa-repeat" title="Republish">&nbsp;</span> Republish</button>
						<button type="button" class="btn btn-link forget"><span class="forget fa fa-eraser" title="Forget">&nbsp;</span> Forget</button>
					</td>
		    	</tr>
		    	{{#Checks}}
		    	<tr
					{{#equal ../Status "1"}}class="warning"{{/equal}}
					{{#equal ../Status "3"}}class="danger"{{/equal}}
		    	 >
					<td></td>
					<td></td>
					<td class="uuid">{{UUID}}</td>
					<td>{{Type}}</td>
					<td>
						<span 
						{{#equal Status "0"}}class="fa fa-ellipsis-h" title="Ignored"{{/equal}}
						{{#equal Status "1"}}class="fa fa-exclamation" title="Invalid"{{/equal}}
						{{#equal Status "2"}}class="fa fa-clock-o" title="In Progress"{{/equal}}
						{{#equal Status "3"}}class="fa fa-times" title="Failed"{{/equal}}
						{{#equal Status "4"}}class="fa fa-check" title="Success"{{/equal}}
						> </span>
					</td>
					<td></td>
		    	</tr>
		    	{{/Checks}}
			</tbody>
				{{/events}}
		</table>
		{{^events}}
		<div>There are no events to display.</div>
		{{/events}}
		<script type="text/javascript">
		$(".forget").on("click", function(e) {
			var row = $(this).parent().parent();
			$.post("/__history/forget", "tid=" + $(".tid", row).text(), function() {
				row.parent().remove();
			});
			return false;
		});
		
		$(".republish").on("click", function(e) {
			$.post("/__history/republish", "uuid=" + $(".uuid", $(this).parent().parent().next()).first().text(), function() {
				alert("Republish job was submitted.");
				location.reload(true);
			});
			return false;
		});
		</script>
		</div>
	</body>
</html>`
)

func init() {
	history = make([]*publishEvent, 0)
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
