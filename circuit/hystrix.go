package circuit

import (
	"context"
	"fmt"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

var (
	COMMAND = "payment-api"
)

func HystrixConfig() {
	hystrix.ConfigureCommand(COMMAND, hystrix.CommandConfig{
		Timeout:                int(5 * time.Second), // 执行command的超时时间
		MaxConcurrentRequests:  5,                    // command的最大并发量
		SleepWindow:            10000,                // 单位是毫秒
		RequestVolumeThreshold: 100,                  // 10秒内请求数量，达到此数后才去判断是否开启熔断
		ErrorPercentThreshold:  20,                   // 错误百分比阈值20%，请求数量大于等于阈值并且错误率达到百分比后启动熔断
	})
}

// clientWrapper 是一个结构体，实现了client.Client接口
type clientWrapper struct {
	client.Client
}

// Call 是clientWrapper结构体的方法，用于进行服务调用
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {

	// 使用hystrix包的Do函数执行服务调用，并设置成功、失败时的处理函数
	return hystrix.Do(COMMAND,
		func() error {
			// 在这里执行实际的服务调用
			return c.Client.Call(ctx, req, rsp, opts...)
		},
		func(err error) error { // 处理调用失败时的降级操作函数
			fmt.Println(err)
			return err
		})
}

// NewClientWrapper 返回一个hystrix的客户端Wrapper
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c} // 返回一个包装后的clientWrapper
	}
}
