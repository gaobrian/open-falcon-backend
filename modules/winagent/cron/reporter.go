package cron

import (
	"fmt"
	"time"

	"github.com/gaobrian/open-falcon-backend/common/model"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/g"
	log "github.com/Sirupsen/logrus"
	"github.com/gaobrian/open-falcon-backend/modules/winagent/tools/os"
)

func ReportAgentStatus() {
	if g.Config().Heartbeat.Enabled && g.Config().Heartbeat.Addr != "" {
		go reportAgentStatus(time.Duration(g.Config().Heartbeat.Interval) * time.Second)
	}
}

func reportAgentStatus(interval time.Duration) {
	for {
		hostname, err := g.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("error:%s", err.Error())
		}

		currPluginVersion, currPluginErr := g.GetCurrPluginVersion()
		if currPluginErr != nil {
			log.Warnln("GetCurrPluginVersion returns: ", currPluginErr)
		}

		currPluginRepo, currRepoErr := g.GetCurrGitRepo()
		if currRepoErr != nil {
			log.Warnln("GetCurrGitRepo returns: ", currRepoErr)
		}


		osinfo, err :=  os.OSVersion()
		if err != nil {
			osinfo = ""
		}

		sysuptime := os.SysupTime2Int64()

		servertime := time.Now().Unix()


		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: currPluginVersion,
			GitRepo:       currPluginRepo,
			Os:            osinfo,
			SysUpTime:     float64(sysuptime),
			ServerTime:    float64(servertime),
			Tags:          g.Config().EndpointTags,
		}

		log.Debugln("show req of Agent.ReportStatus: ", req)
		var resp model.SimpleRpcResponse
		err = g.HbsClient.Call("Agent.ReportStatus", req, &resp)
		if err != nil || resp.Code != 0 {
			log.Errorln("call Agent.ReportStatus fail:", err, "Request:", req, "Response:", resp)
		}

		time.Sleep(interval)
	}
}
