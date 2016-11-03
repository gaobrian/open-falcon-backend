package funcs

import (
	"github.com/gaobrian/open-falcon-backend/common/model"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/nux"
)

func SocketStatSummaryMetrics() (L []*model.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Println(err)
		return
	}

	for k, v := range ssMap {
		L = append(L, GaugeValue("ss."+k, v))
	}

	return
}