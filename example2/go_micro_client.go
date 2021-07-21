package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/web"
	"github.com/gin-gonic/gin"
	"protos/Services"
	"time"
)
//go get -u github.com/micro/go-grpc

type logWrapper struct {
      client.Client
}
func (this *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+"调用接口"+"[wrapper] client request service: %s method: %s\n", req.Service(), req.Endpoint())
	return this.Client.Call(ctx, req, rsp)
}

func NewlogWrap(c client.Client) client.Client {
	return &logWrapper{c}
}
func main(){
	// Register consul
	consulReg := consul.NewRegistry(registry.Addrs(":8500"))
	//gin 路由
	ginRouter:=gin.New()

	httpServer:=web.NewService(
		web.Name("httpprodservice"),
		web.Address(":8099"),
		web.Handler(ginRouter),
		web.Registry(consulReg),

		// 为注册的服务添加Metadata，指定请求协议为http
		web.Metadata(map[string]string{"protocol" : "http"}),
	)
	myService:=micro.NewService(micro.Name("prodservice.client"),
		micro.Registry(consulReg),
		micro.WrapClient(NewlogWrap),
		)
	myService.Init()
	prodService := Services.NewProdService("prodservice",myService.Client())

	v1group:=ginRouter.Group("/v1")
	{
		v1group.Handle("POST","/prods", func(gincontext *gin.Context) {
			var prodRes Services.ProdsRequest
			err :=gincontext.Bind(&prodRes)
			if err!=nil{
				gincontext.JSON(500,gin.H{"status":err.Error()})

			}else {
				prodRes,_:=prodService.GetProdsList(context.Background(),&prodRes)
				gincontext.JSON(200,gin.H{"data":prodRes.Data})
			}

		})

	}

	httpServer.Init()
	httpServer.Run()
}
