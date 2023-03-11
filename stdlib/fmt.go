package stdlib

import (
	"time"
)

func Print(msg string) {
	println(msg)
}

func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
