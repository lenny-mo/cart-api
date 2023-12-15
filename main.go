package main

import (
	"fmt"
	"net/http"

	"github.com/lenny-mo/cart-api/handler"
	"github.com/lenny-mo/cart-api/proto/cartapi"
	"github.com/lenny-mo/cart-api/utils"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/lenny-mo/cart/proto/cart"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {

	// 2 注册中心 使用consul v2
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	// 3 创建链路追踪
	tracer, tracerio, err := utils.NewTracer("cart-server-api", "127.0.0.1:6831")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer tracerio.Close()
	opentracing.SetGlobalTracer(tracer) // 设置全局的链路追踪

	// 4 client加载熔断器
	// StreamHandler 是一个 HTTP 处理器（handler），它用于提供一个实时的数据流，这个数据流包含了 hystrix 熔断器的度量信息和状态。
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()

	// 启动一个 HTTP 服务器
	// 本地的 8090 端口上启动一个新的服务器实例，并且所有到这个端口的 HTTP 请求都应该由 hystrixStreamHandler 来处理。
	go func() {
		err := http.ListenAndServe("127.0.0.1:8090", hystrixStreamHandler)
		if err != nil {
			fmt.Println(err)
		}
	}()

	service := micro.NewService(
		micro.Name("go.micro.api.github.com/lenny-mo/cart-api"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8090"),
		micro.Registry(consulRegistry),
		// 充当client访问其他服务时，加载链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(utils.NewClientHystrixWrapper()),
	)
	// 创建轮询 RoundRobin 负载均衡策略
	rrSelector := selector.NewSelector(
		selector.SetStrategy(selector.RoundRobin),
	)
	// 在client 添加负载均衡 轮巡算法
	service.Init(micro.Selector(rrSelector))

	cartService := cart.NewCartService("go.micro.service.cart", service.Client())

	// 注册服务
	if err := cartapi.RegisterCartApiHandler(service.Server(), &handler.CartAPI{
		CarService: cartService,
	}); err != nil {
		fmt.Println("error when register service!")
	}

	// hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
	// 	Timeout:                1000, // 超时时间设置为1000毫秒
	// 	MaxConcurrentRequests:  1000, // 最大并发请求数
	// 	RequestVolumeThreshold: 20,   // 请求量阈值
	// 	SleepWindow:            5000, // 熔断器打开后，尝试再次关闭之前的等待时间 /ms
	// 	ErrorPercentThreshold:  50,   // 触发熔断的错误百分比阈值
	// })

	// 启动服务
	if err := service.Run(); err != nil {
		fmt.Println("error during starting the service")
	}
}
