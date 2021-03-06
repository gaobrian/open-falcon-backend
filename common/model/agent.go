package model

import (
	"fmt"
)

type AgentReportRequest struct {
	Hostname      string
	IP            string
	AgentVersion  string
	PluginVersion string
	GitRepo       string
	//Add more fields for system
	Os            string
	SysUpTime     float64
	ServerTime    float64
        Tags          string
}

func (this *AgentReportRequest) String() string {
	return fmt.Sprintf(
		"<Hostname:%s, IP:%s, AgentVersion:%s, PluginVersion:%s, GitRepo: %s>",
		this.Hostname,
		this.IP,
		this.AgentVersion,
		this.PluginVersion,
		this.GitRepo,
	)
}

type AgentUpdateInfo struct {
	LastUpdate    int64
	ReportRequest *AgentReportRequest
}

type AgentHeartbeatRequest struct {
	Hostname string
	Checksum string
}

func (this *AgentHeartbeatRequest) String() string {
	return fmt.Sprintf(
		"<Hostname: %s, Checksum: %s>",
		this.Hostname,
		this.Checksum,
	)
}

type AgentPluginsResponse struct {
	Plugins       []string
	Timestamp     int64
	GitRepo       string
	GitUpdate     bool
	GitRepoUpdate bool
}

func (this *AgentPluginsResponse) String() string {
	return fmt.Sprintf(
		"<Plugins:%v, Timestamp:%v, GitRepo:%v, GitUpdate:%v, GitRepoUpdate:%v>",
		this.Plugins,
		this.Timestamp,
		this.GitRepo,
		this.GitUpdate,
		this.GitRepoUpdate,
	)
}

// e.g. net.port.listen or proc.num
type BuiltinMetric struct {
	Metric string
	Tags   string
}

func (this *BuiltinMetric) String() string {
	return fmt.Sprintf(
		"%s/%s",
		this.Metric,
		this.Tags,
	)
}

type BuiltinMetricResponse struct {
	Metrics   []*BuiltinMetric
	Checksum  string
	Timestamp int64
}

func (this *BuiltinMetricResponse) String() string {
	return fmt.Sprintf(
		"<Metrics:%v, Checksum:%s, Timestamp:%v>",
		this.Metrics,
		this.Checksum,
		this.Timestamp,
	)
}

type BuiltinMetricSlice []*BuiltinMetric

func (this BuiltinMetricSlice) Len() int {
	return len(this)
}
func (this BuiltinMetricSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this BuiltinMetricSlice) Less(i, j int) bool {
	return this[i].String() < this[j].String()
}
