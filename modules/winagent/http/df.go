package http

import (
	"github.com/gaobrian/open-falcon-backend/modules/winagent/funcs"
	"net/http"
)

func configDfRoutes() {
	http.HandleFunc("/page/df", func(w http.ResponseWriter, r *http.Request) {
		usage := funcs.DeviceMetrics()
		RenderDataJson(w, usage)

	})
}
