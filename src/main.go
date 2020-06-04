package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	"github.com/linclin/gopub/src/models"
	"github.com/linclin/gopub/src/routers"
	"github.com/linclin/gopub/src/util"
	"os"
)

func initArgs() {
	args := os.Args
	for _, v := range args {
		if v == "-syncdb" {
			models.Syncdb()
			os.Exit(0)
		}
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error("Panic error:", err)
		}
	}()
	initArgs()
	logs.Info(beego.BConfig.RunMode)
	util.Graceful()
	util.Swagger()
	toolbox.StartTask()
	// init_sever.Start()
	routers.Run()
	beego.Run()
}
