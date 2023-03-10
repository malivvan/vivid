package vm

import "time"

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func print(msg string) {
	println(msg)
}
