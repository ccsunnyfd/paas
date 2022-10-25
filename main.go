package main

import (
	"fmt"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"saas/governance/database"
	"saas/governance/trace"

	config "saas/governance/config"
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

	// micro服务
	service := micro.NewService(
		micro.Name("paas"),
		micro.Version("1.0"),
		// 注册中心
		micro.Registry(registryCenter),
		// 链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
	)

	service.Init()

	service.Run()
}
