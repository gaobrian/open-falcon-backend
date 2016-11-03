package funcs

import (
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/load"
	"log"
	"github.com/gaobrian/open-falcon-backend/common/model"
)

func LoadMetrics() (L []*model.MetricValue) {

	loadVal, err := load.LoadAvg()
	if err != nil {
		log.Println("Get load fail: ", err)
		return nil
	}

	L = append(L, CounterValue("load.load1", loadVal.Load1))
	L = append(L, CounterValue("load.load5", loadVal.Load5))
	L = append(L, CounterValue("load.load15", loadVal.Load15))

	return
}
