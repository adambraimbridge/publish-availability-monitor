package history

import (
	"log"
	"os"
)

const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var infoLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "INFO  - ", logPattern)
}
