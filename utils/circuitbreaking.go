package utils

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type ClientWrapper struct {
	client.Client
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &ClientWrapper{c}
	}
}

// Call 是 clientWrapper 结构体的一个方法，它实现了一个熔断器模式。
// 这个方法使用了 Hystrix 库来包装原始的 RPC 调用，以提供熔断和降级功能。
func (c *ClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// hystrix.Do 方法用于执行包含熔断逻辑的函数。
	// 第一个参数是一个唯一的命令名称，通常由服务名和端点名组合而成。
	// 这个名称用于在熔断器的内部逻辑中标识和区分不同的命令。
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 这是正常执行的路径，即当熔断器处于关闭状态时执行的代码。
		// 这里打印出服务名和端点名，用于调试或日志记录。
		fmt.Println(req.Service() + "." + req.Endpoint())

		// 调用原始的客户端来执行实际的 RPC 调用。
		// 这里传递了上下文（ctx）、请求对象（req）和响应对象（rsp）。
		// opts... 是可变参数，表示调用选项。
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		// 这是熔断器打开时执行的路径，即当服务不可用或响应时间过长时执行的代码。
		// 这里打印出错误信息，用于调试或日志记录。
		fmt.Println(err)

		// 返回错误，通常这里可以实现降级逻辑，即当服务不可用时的备选方案。
		return err
	})
}
