package funcs

import (
	"github.com/gaobrian/open-falcon-backend/common/model"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/g"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

const (
	IntervalThirtySec = 30
)

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := g.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				AgentMetricsThirty,
			},
			Interval: IntervalThirtySec,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				AgentMetrics,
				CpuMetrics,
				NetMetrics,
				MemMetrics,
				DiskIOMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				UrlMetrics,
			},
			Interval: interval,
		},
	}
}
