package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type singValues struct {
	isDebug bool
	isLogOn bool
	once    sync.Once
}

var vals singValues

func getSingValues() singValues {
	vals.once.Do(func() {
		vals.isDebug = strings.ToLower(os.Getenv("GO_ENV")) != "prod"
		vals.isLogOn, _ = strconv.ParseBool(os.Getenv("LOG_ON"))

		if !vals.isLogOn {
			fmt.Println("(env LOG_ON=false) logs disabled")
		}
	})

	return vals
}

func getFormattedMessage(msg string) string {
	return fmt.Sprintf("%s - %s", time.Now().Format("02/01/2006 15:04:05"), msg)
}

func Info(msg string, args ...interface{}) {
	if getSingValues().isLogOn {
		fmtMsg := getFormattedMessage(msg)
		if getSingValues().isDebug && len(args) > 0 {
			fmt.Println(fmtMsg, args)
		} else {
			fmt.Println(fmtMsg)
		}
	}
}

func Error(msg string, errs ...error) {
	Info(msg, errs)
}
