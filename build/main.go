package main

import (
	"log"

	"gateway.oscto.icu/router"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/web"
)

var (
	GatewayServiceName    = "gateway.oscto.icu"
	GatewayServiceAddress = ":8080"
)

func main() {

	service := web.NewService(
		web.Name(GatewayServiceName),
		web.Address(GatewayServiceAddress),
	)
	_ = service.Init()
	r := gin.Default()
	router.Account(r)
	router.Location(r)
	router.Image(r)
	service.Handle("/", r)
	if err := service.Run(); err != nil {
		log.Panicln("启动错误,", err)
	}
}
