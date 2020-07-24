package program

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"runtime"
	"runtime/debug"
	"time"
)

func WithRecover(handler func()) {
	//todo we can handle program num here
	go func() {
		defer DumpStack(false)
		handler()
	}()
}

func DumpStack(exitIFPanic bool) {
	if err := recover(); err != nil {
		fmt.Println("program num ", runtime.NumGoroutine())
		logrus.WithField("obj", err).Error("Fatal error occurred.")
		var buf bytes.Buffer
		stack := debug.Stack()
		buf.WriteString(fmt.Sprintf("Panic: %v\n", err))
		buf.Write(stack)
		dumpName := "dump_" + time.Now().Format("20060102-150405")
		nerr := ioutil.WriteFile(dumpName, buf.Bytes(), 0644)
		if nerr != nil {
			fmt.Println("write dump file error", nerr)
			fmt.Println(buf.String())
		}

		if exitIFPanic {
			logrus.Fatalf("panic %v ", buf.String())
		}
	}
}
