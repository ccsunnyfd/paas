package main

import (
	"fmt"
	ratelimiter "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"saas/governance/breaker"
	"saas/governance/database"
	"saas/governance/monitor"
	"saas/governance/trace"

	"saas/governance/config"
	"saas/governance/registry"
)

func main() {
	// 注册中心
	registryCenter := registry.GetRegistry("localhost", 8500)

	// 配置中心
	configCenter, err := config.GetConfig("localhost", 8500, "/micro/config")
	if err != nil {
		panic(err)
	}

	// 数据库
	mysqlConfig := config.GetMysqlFromConfig(configCenter, "mysql")

	_, err = database.GetMysql(mysqlConfig.User, mysqlConfig.Pwd, mysqlConfig.Database)
	if err != nil {
		panic(err)
	}

	fmt.Println("连接mysql成功")

	// 链路追踪
	t, io, err := trace.NewTracer("paas", "localhost", 6831)
	if err != nil {
		panic(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断器
	breaker.StartListening("0.0.0.0", "9092")

	// 监控
	monitor.PrometheusBoot(9192)

	// micro服务
	service := micro.NewService(
		micro.Name("paas"),
		micro.Version("1.0"),
		// 注册中心
		micro.Registry(registryCenter),
		// 链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 熔断，只作为客户端的时候起作用
		micro.WrapClient(breaker.NewClientHystrixWrapper()),
		// 限流，作为服务端的时候起作用
		micro.WrapHandler(ratelimiter.NewHandlerWrapper(1000)),
	)

	service.Init()

	service.Run()
}
