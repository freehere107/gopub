package util

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(c chan os.Signal) {
	switch <-c {
	case syscall.SIGINT, syscall.SIGTERM:

		logs.Info("Shutdown quickly, bye...")
	case syscall.SIGQUIT:
		logs.Info("Shutdown gracefully, bye...")
		// do graceful shutdown
	}
	os.Exit(0)
}

func Graceful() {
	graceful, _ := beego.AppConfig.Bool("Graceful")
	if graceful {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go handleSignals(sigs)
	}
}
