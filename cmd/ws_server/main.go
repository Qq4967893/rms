package main

import (
	"backend-rms/bootstrap"
	"backend-rms/ws"
	"backend-rms/ws/middleware"
	"backend-rms/ws/route"
	"time"
)

func main() {
	time.Local = time.UTC

	bootstrap.InitConfig()
	bootstrap.InitLogger()
	bootstrap.InitJWT()
	bootstrap.InitRedis()
	bootstrap.InitDB()

	router := ws.NewEngine()
	router.Use(middleware.Authenticate())
	route.RegisterRoutes(router)
	err := router.Run(":8889")
	if err != nil {
		panic(err)
	}
}
