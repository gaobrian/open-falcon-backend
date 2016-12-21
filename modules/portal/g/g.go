package g

import (
	"runtime"
)

const (
	VERSION = "1.0.5"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
