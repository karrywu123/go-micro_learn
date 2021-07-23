package main

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/web"
	"go-micro_learn/example3/protos/Services"
	"go-micro_learn/example3/router"
	go_micro_wappers "go-micro_learn/example3/go-micro-wappers"
	"strconv"
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

func NewProd(id int32,pname string) *Services.ProdModel {
	return &Services.ProdModel{ProdID:id,ProdName:pname}
}

func defaultProds()(*Services.ProdListResponse,error){
	models:=make([]*Services.ProdModel,0)
	var i int32
	for i=0;i<2.;i++{
		models=append(models,NewProd(20+i,"prodname"+strconv.Itoa(20+int(i))))
	}
	res:=&Services.ProdListResponse{}
	res.Data=models
	return res,nil
}

func main(){
	// Register consul
	consulReg := consul.NewRegistry(registry.Addrs(":8500"))
	myService:=micro.NewService(micro.Name("prodservice.client"),
		micro.Registry(consulReg),
		micro.WrapClient(NewlogWrap),//记录访问日志
		//加熔断 超时熔断
		micro.WrapClient(go_micro_wappers.NewlogWrap),

	)
	myService.Init()
	prodService := Services.NewProdService("prodservice",myService.Client())

	httpServer:=web.NewService(
		web.Name("httpprodservice"),
		web.Address(":8099"),
		web.Handler(router.NewGinRouter(prodService)),
		web.Registry(consulReg),

		// 为注册的服务添加Metadata，指定请求协议为http
		web.Metadata(map[string]string{"protocol" : "http"}),
	)


	httpServer.Init()
	httpServer.Run()
}
