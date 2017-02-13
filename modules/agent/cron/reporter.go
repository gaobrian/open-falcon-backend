package cron

import (
	"fmt"
	"time"

	"github.com/gaobrian/open-falcon-backend/common/model"
	"github.com/gaobrian/open-falcon-backend/modules/agent/g"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/sys"
	"github.com/toolkits/file"
	"strings"
	"strconv"
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

		osinfo, err :=  sys.CmdOutNoLn("uname", "-r")
		if err != nil {
			osinfo = ""
		}

		content, err := file.ToTrimString("/proc/uptime")
		if err != nil {
			content = "0 0"
		}

		fields := strings.Fields(content)
		secStr := fields[0]
		var sysuptime float64
		sysuptime, err = strconv.ParseFloat(secStr, 64)
		if err != nil {
			sysuptime = 0
		}

                servertime := time.Now().Unix()

		req := model.AgentReportRequest{
			Hostname:      hostname,
			IP:            g.IP(),
			AgentVersion:  g.VERSION,
			PluginVersion: currPluginVersion,
			GitRepo:       currPluginRepo,
			Os:            osinfo,
			SysUpTime:     sysuptime,
			ServerTime:    servertime,
			Tags:          "",
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
