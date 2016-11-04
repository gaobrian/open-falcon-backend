package http

import (
	"github.com/gaobrian/open-falcon-backend/modules/winagent/funcs"
	"net/http"
)

func configIoStatRoutes() {
	http.HandleFunc("/page/diskio", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.DiskIOMetrics())
	})
}
