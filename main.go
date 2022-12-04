package main

import (
	"context"
	pb "git.vonechain.com/vone-bfs/gateway/proto"
	"log"
	"sync"

	"git.vonechain.com/vone-bfs/gateway/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/cli/debug/trace/jaeger"
	grpcc "github.com/go-micro/plugins/v4/client/grpc"
	httpServer "github.com/go-micro/plugins/v4/server/http"
	ot "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
)

var (
	service     = "gateway"
	SERVER_NAME = "gateway-http"
	version     = "latest"
)

func main() {

	httpSrv := httpServer.NewServer(
		server.Name(SERVER_NAME),
		server.Address(":8090"),
	)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// register router
	demo := handler.NewDemo()
	demo.InitRouter(router)

	hd := httpSrv.NewHandler(router)
	if err := httpSrv.Handle(hd); err != nil {
		log.Fatalln(err)
	}

	// Create tracer
	tracer, closer, err := jaeger.NewTracer(
		jaeger.Name(service),
		jaeger.FromEnv(true),
		jaeger.GlobalTracer(true),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer closer.Close()

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// Create service
	srv := micro.NewService(
		micro.Server(httpSrv),
		//micro.Server(grpcs.NewServer()),
		micro.Client(grpcc.NewClient()),
		micro.BeforeStart(func() error {
			logger.Infof("Starting service %s", service)
			return nil
		}),
		micro.BeforeStop(func() error {
			logger.Infof("Shutting down service %s", service)
			cancel()
			return nil
		}),
		micro.AfterStop(func() error {
			wg.Wait()
			return nil
		}),
		micro.WrapCall(ot.NewCallWrapper(tracer)),
		micro.WrapClient(ot.NewClientWrapper(tracer)),
		micro.WrapHandler(ot.NewHandlerWrapper(tracer)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(tracer)),
	)
	srv.Init(
		micro.Name(service),
		micro.Version(version),
	)
	srv.Server().Init(
		server.Wait(&wg),
	)

	ctx = server.NewContext(ctx, srv.Server())

	// Register handler
	if err := pb.RegisterGatewayHandler(srv.Server(), new(handler.Gateway)); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
