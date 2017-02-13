package http

import (
	"github.com/gaobrian/open-falcon-backend/modules/winagent/g"
	"net/http"
	"runtime"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/os"
)

func configKernelRoutes() {
	http.HandleFunc("/proc/kernel/hostname", func(w http.ResponseWriter, r *http.Request) {
		data, err := g.Hostname()
		AutoRender(w, data, err)
	})

	http.HandleFunc("/proc/kernel/maxproc", func(w http.ResponseWriter, r *http.Request) {
		data := runtime.NumCPU()
		AutoRender(w, data, nil)
	})

	http.HandleFunc("/proc/kernel/version", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.OSVersion()
		AutoRender(w, data, err)
	})

}
