package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gaobrian/open-falcon-backend/common/logruslog"
	"github.com/gaobrian/open-falcon-backend/common/vipercfg"

	"github.com/gaobrian/open-falcon-backend/modules/query/conf"
	"github.com/gaobrian/open-falcon-backend/modules/query/database"
	"github.com/gaobrian/open-falcon-backend/modules/query/g"
	ginHttp "github.com/gaobrian/open-falcon-backend/modules/query/gin_http"
	"github.com/gaobrian/open-falcon-backend/modules/query/graph"
	"github.com/gaobrian/open-falcon-backend/modules/query/grpc"
	"github.com/gaobrian/open-falcon-backend/modules/query/http"
	"github.com/gaobrian/open-falcon-backend/modules/query/proc"
)

func main() {
	vipercfg.Parse()
	vipercfg.Bind()

	if vipercfg.Config().GetBool("version") {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	// config
	vipercfg.Load()
	g.ParseConfig(vipercfg.Config().GetString("config"))
	logruslog.Init()
	gconf := g.Config()
	// proc
	proc.Start()

	// graph
	go graph.Start()

	if gconf.Grpc.Enabled {
		// grpc
		go grpc.Start()
	}

	if gconf.GinHttp.Enabled {
		//lambdaSetup
		database.Init()
		conf.ReadConf()
		go ginHttp.StartWeb()
	}

	if gconf.Http.Enabled {
		// http
		go http.Start()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		if sig.String() == "^C" {
			os.Exit(3)
		}
	}
}
