package logformat

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	timeFmt = "2006-01-02 15:04:05.000"
	txIdKey = "transaction_id"
)

// Our default SLF4J format is: "%-5p [%d{ISO8601, GMT}] %c: %X{transaction_id} %m [%thread]%n%xEx"
type SLF4JFormatter struct {
	stackPattern *regexp.Regexp
}

func NewSLF4JFormatter(pattern string) *SLF4JFormatter {
	var p *regexp.Regexp
	if pattern != "" {
		p = regexp.MustCompile(pattern)
	}

	f := SLF4JFormatter{p}
	return &f
}

func (f *SLF4JFormatter) Format(entry *log.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	// except for warn, which by default becomes WARNING
	if entry.Level == log.WarnLevel {
		level = "WARN"
	}

	timestamp := strings.Replace(entry.Time.UTC().Format(timeFmt), ".", ",", 1)
	codeLocation := f.findCodeLocation()
	tx := f.findTransactionId(entry.Data)

	msg := fmt.Sprintf("%-5s [%s] %s %s %s\n", level, timestamp, codeLocation, tx, entry.Message)

	return []byte(msg), nil
}

func (f *SLF4JFormatter) findCodeLocation() string {
	if f.stackPattern == nil {
		return ""
	}

	// start at 2 because we know 0 and 1 are within this file
	for i := 2; i < 10; i++ {
		_, file, lineNum, ok := runtime.Caller(i)
		//		os.Stdout.WriteString("checking for pattern " + f.stackPattern.String() + " in stack " + file)
		if ok && f.stackPattern.MatchString(file) {
			return fmt.Sprintf("%s:%v:", file[strings.LastIndex(file, "/")+1:], lineNum)
		}
	}

	return "unknown"
}

func (f *SLF4JFormatter) findTransactionId(data map[string]interface{}) string {
	var tx string
	if v, found := data[txIdKey]; found {
		tx = v.(string)
		if tx != "" && !strings.HasPrefix(tx, txIdKey) {
			tx = txIdKey + "=" + tx
		}
	}

	return tx
}
