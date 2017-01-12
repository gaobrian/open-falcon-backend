package funcs

import (
	"log"

	"github.com/gaobrian/open-falcon-backend/common/model"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/net"
)


func TcpipMetrics() (L []*model.MetricValue) {

	ds, err := net.TcpipCounters()

	if err != nil {
		log.Println("Get tcpip data fail: ", err)
		return
	}

	L = append(L, CounterValue("tcpip.confailures", ds[0].ConFailures))
	L = append(L, CounterValue("tcpip.conactive", ds[0].ConActive))
	L = append(L, CounterValue("tcpip.conpassive", ds[0].ConPassive))
	L = append(L, GaugeValue("tcpip.conestablished", ds[0].ConEstablished))
	L = append(L, CounterValue("tcpip.conreset", ds[0].ConReset))

	return
}
