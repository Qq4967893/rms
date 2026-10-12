package main

import (
	"backend-rms/bootstrap"
	"backend-rms/core/logger"
	"backend-rms/cron/calc_profit"
	"backend-rms/cron/calc_spread_profit"
	"backend-rms/grpc"
	"context"
	"net/http"
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
	bootstrap.InitSentry()
	bootstrap.InitJWT()
	bootstrap.InitDB()
	bootstrap.InitRedis()
	bootstrap.InitMongoDB()
	bootstrap.InitCasbin()
	grpc.InitTSPClient()

	go calc_profit.Start()
	go calc_spread_profit.Start()

	r := bootstrap.InitRouter()
	srv := &http.Server{
		Addr:    ":8888",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorIf(err)
			quit <- syscall.SIGQUIT
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.ErrorIf(srv.Shutdown(ctx))
}
