package home

import (
	"github.com/gaobrian/open-falcon-backend/modules/portal/g"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/gaobrian/open-falcon-backend/modules/portal/controllers/base"
)

func init()  {
	ConfigRoutes()
}

func ConfigRoutes() {
	beego.InsertFilter("/",beego.BeforeRouter,base.FilterLoginUser)
	beego.Router("/", &HomeController{})

	beego.Get("/health", func(ctx *context.Context) {
		ctx.Output.Body([]byte("ok"))
	})

	beego.Get("/version", func(ctx *context.Context) {
		ctx.Output.Body([]byte(g.VERSION))
	})
}
