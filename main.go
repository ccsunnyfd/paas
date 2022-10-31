package main

import (
	"flag"
	"fmt"
	ratelimiter "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/mitchellh/go-homedir"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/util/log"
	"path/filepath"
	v1 "saas/api/pod/v1"
	"saas/governance/breaker"
	"saas/governance/monitor"
	"saas/governance/trace"
	"saas/internal/biz"
	"saas/internal/data"
	service2 "saas/internal/service"

	"saas/governance/config"
	"saas/governance/registry"
)

func main() {
	// 注册中心
	registryCenter := registry.GetRegistry("localhost", 8500)

	// 配置中心
	configCenter, err := config.GetConfig("localhost", 8500, "/micro/config")
	if err != nil {
		log.Fatal(err)
	}

	// 数据库
	mysqlConfig := config.GetMysqlFromConfig(configCenter, "mysql")

	db, err := data.NewMysql(mysqlConfig.User, mysqlConfig.Pwd, mysqlConfig.Database)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("连接mysql成功")

	// 链路追踪
	t, io, err := trace.NewTracer("paas", "localhost", 6831)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	// 熔断器
	breaker.StartListening("0.0.0.0", "9092")

	// 监控
	monitor.PrometheusBoot(9192)

	// K8S连接
	var kubeconfig *string
	if home, err := homedir.Dir(); err != nil {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig file 在当前系统中的地址")
	} else {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig file 在当前系统中的地址")
	}
	flag.Parse()

	k8sClientset, err := data.NewK8SClientset(*kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// micro服务
	service := micro.NewService(
		micro.Name("paas.service.pod"),
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

	podRepo := data.NewPodRepo(db)
	err = podRepo.InitTable()
	if err != nil {
		log.Fatal(err)
	}

	k8sRepo := data.NewK8SRepo(k8sClientset)

	podUsercase := biz.NewPodUsecase(podRepo, k8sRepo)
	podHandler := service2.NewPodHandler(podUsercase)
	err = v1.RegisterPodHandler(service.Server(), podHandler)
	if err != nil {
		log.Fatal(err)
	}

	if err = service.Run(); err != nil {
		log.Fatal(err)
	}
}
