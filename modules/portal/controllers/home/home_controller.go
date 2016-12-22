package home

import (
	"github.com/gaobrian/open-falcon-backend/modules/portal/g"
        _ "github.com/gaobrian/open-falcon-backend/modules/portal/models/uic"
	"github.com/astaxie/beego"
	_ "log"
)

type HomeController struct {
	beego.Controller
}

func (this *HomeController) Get() {
	this.Data["Shortcut"] = g.Config().Shortcut
	this.TplName = "home/index.html"
}
