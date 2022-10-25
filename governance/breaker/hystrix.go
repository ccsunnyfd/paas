package breaker

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"go-micro.dev/v4/client"
	"net"
	"net/http"
)

func StartListening(host, port string) {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 启动监听程序
	go func() {
		// http://localhost:9092/turbine/turbine.stream
		// 看板访问地址 http://127.0.0.1:9002/hystrix
		err := http.ListenAndServe(net.JoinHostPort(host, port),
			hystrixStreamHandler)
		if err != nil {
			fmt.Println(err)
		}
	}()
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 正常执行
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		// 走熔断逻辑，每个服务都不一样
		fmt.Println(req.Service() + "." + req.Endpoint() + "的熔断逻辑")
		return e
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
