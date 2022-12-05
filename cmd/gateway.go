package main

import (
	"context"
	"fmt"
	"git.vonechain.com/vone-bfs/gateway/handler"
	"git.vonechain.com/vone-bfs/gateway/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/cli/debug/trace/jaeger"
	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/etcd"
	httpServer "github.com/go-micro/plugins/v4/server/http"
	ot "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"net/http"
	"sync"
	"time"
)

var (
	NAME    = "gateway"
	ADDR    = "0.0.0.0:8090"
	VERSION = "0.1"
)

func main() {
	// 初始化配置
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	// 跟踪初始化
	tracer, closer, err := jaeger.NewTracer(
		jaeger.Name(NAME),
		jaeger.FromEnv(true),
		jaeger.GlobalTracer(true),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer closer.Close()

	// 服务初始化
	srv := micro.NewService(
		micro.Server(httpServer.NewServer()),
		micro.Client(grpc.NewClient()),
		// 启动
		micro.BeforeStart(func() error {
			logger.Infof("Starting service %s", NAME)
			return nil
		}),
		micro.BeforeStop(func() error {
			logger.Infof("Shutting down service %s", NAME)
			cancel()
			return nil
		}),
		micro.AfterStop(func() error {
			wg.Wait()
			return nil
		}),

		// 服务注册
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(30*time.Second),
		micro.Registry(etcd.NewRegistry(registry.Addrs("43.138.199.52:2379"))),
		// 跟踪
		micro.WrapCall(ot.NewCallWrapper(tracer)),
		micro.WrapClient(ot.NewClientWrapper(tracer)),
		micro.WrapHandler(ot.NewHandlerWrapper(tracer)),
		micro.WrapSubscriber(ot.NewSubscriberWrapper(tracer)),
	)

	opt := []micro.Option{
		micro.Name(NAME),
		micro.Address(ADDR),
		micro.Version(VERSION),
	}

	srv.Init(opt...)
	srv.Server().Init(server.Wait(&wg))
	srvs := handler.NewServices(srv.Client())
	server.NewContext(ctx, srv.Server())
	httpSrv := gin.Default()
	httpSrv.Use(middleware.Casbin())
	httpSrv.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "success", "controller": "index"})
	})

	httpSrv.GET("/user", srvs.Register)

	if err := micro.RegisterHandler(srv.Server(), httpSrv); err != nil {
		fmt.Println(fmt.Errorf("micro register handler error: %v", err))
	}

	if err := srv.Run(); err != nil {
		fmt.Println(fmt.Errorf("failed to service start, error: %w", err))
	}
}
