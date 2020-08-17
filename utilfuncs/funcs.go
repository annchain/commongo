package utilfuncs

import (
	"fmt"
	"os"
	"time"
)

func PanicIfError(err error, message string) {
	if err != nil {
		fmt.Println("panic: " + message)
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func DrainTimer(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func DrainTicker(timer *time.Ticker) {
	select {
	case <-timer.C:
	default:
	}
}
