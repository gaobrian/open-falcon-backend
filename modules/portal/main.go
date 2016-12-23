package main

import (
	_ "github.com/gaobrian/open-falcon-backend/modules/portal/routers"
	"github.com/astaxie/beego"
	"fmt"
	"github.com/gaobrian/open-falcon-backend/modules/portal/cache"
	"github.com/gaobrian/open-falcon-backend/common/logruslog"
	"github.com/gaobrian/open-falcon-backend/modules/portal/g"
	"github.com/gaobrian/open-falcon-backend/common/vipercfg"
	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/logger"
	"os"
	"github.com/gaobrian/open-falcon-backend/modules/portal/models"
	"github.com/gaobrian/open-falcon-backend/modules/portal/models/uic"
	"strings"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/logs"
)

func main() {

	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// parse config
	vipercfg.Load()
	if err := g.ParseConfig(vipercfg.Config().GetString("config")); err != nil {
		log.Fatalln(err)
	}
	logruslog.Init()

	logger.SetLevelWithDefault(g.Config().Log, "info")

	cache.InitCache()

	model.InitDatabase()

	switch strings.ToLower(g.Config().Log) {
	case "info":
		beego.SetLevel(beego.LevelInformational)
		logs.EnableFuncCallDepth(true)
	case "debug":
		beego.SetLevel(beego.LevelDebug)
		logs.EnableFuncCallDepth(true)
	case "warn":
		beego.SetLevel(beego.LevelWarning)
		logs.EnableFuncCallDepth(true)
	case "error":
		beego.SetLevel(beego.LevelError)
		logs.EnableFuncCallDepth(true)
	}


	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))

	beego.AddFuncMap("member", uic.MembersByTeamId)


	beego.SetStaticPath("assets","static/assets")
	beego.Run()
}

