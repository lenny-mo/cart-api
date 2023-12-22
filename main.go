package main

import (
	"fmt"

	"github.com/lenny-mo/cart-api/circuit"
	"github.com/lenny-mo/cart-api/handler"
	"github.com/lenny-mo/cart-api/proto/cartapi"
	"github.com/lenny-mo/emall-utils/tracer"

	"github.com/lenny-mo/cart/proto/cart"
	"github.com/micro/go-micro/v2"
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

	serviceName := "go.micro.api.cart-api"
	// 3 链路追踪
	err := tracer.InitTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tracer.Closer.Close()
	opentracing.SetGlobalTracer(tracer.Tracer)

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8088"),
		micro.Registry(consulRegistry),
		// 充当client访问其他服务时，加载链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		// 添加熔断
		micro.WrapClient(circuit.NewClientWrapper()),
	)

	service.Init()

	cartService := cart.NewCartService("go.micro.service.cart", service.Client())

	// 注册服务
	if err := cartapi.RegisterCartApiHandler(service.Server(), &handler.CartAPI{
		CarService: cartService,
	}); err != nil {
		fmt.Println("error when register service!")
	}

	// 启动服务
	if err := service.Run(); err != nil {
		fmt.Println("error during starting the service")
	}
}
