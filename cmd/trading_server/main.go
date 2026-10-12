package main

import (
	"backend-rms/bootstrap"
	"backend-rms/core/kafka"
	"backend-rms/cron/calc_exposure_usd"
	"backend-rms/tr/init_data"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//设置当前系统时区为UTC时区
	time.Local = time.UTC

	bootstrap.InitConfig()
	bootstrap.InitLogger()
	bootstrap.InitDB()
	bootstrap.InitRedis()
	init_data.InitData()
	init_data.InitGroupData()

	go kafka.Start()
	go calc_exposure_usd.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		s := <-quit
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("admin exit: " + s.String())
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
