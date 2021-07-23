package go_micro_wappers

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/asim/go-micro/v3/client"
	"go-micro_learn/example3/protos/Services"
	"strconv"
	"time"
)


type ProdsWrapper struct {
	client.Client
}

func NewProd(id int32,pname string) *Services.ProdModel {
	return &Services.ProdModel{ProdID:id,ProdName:pname}
}


func defaultData(rsp interface{}){
	switch t:=rsp.(type){
	case *Services.ProdListResponse:
		defaultProds(rsp)
	case *Services.ProdDetailResponse:
		t.Data=NewProd(10,"商品降级 ")
	}
}


func defaultProds(rsp interface{}){
	models:=make([]*Services.ProdModel,0)
	var i int32
	for i=0;i<2.;i++{
		models=append(models,NewProd(20+i,"prodname"+strconv.Itoa(20+int(i))))
	}
	result:=rsp.(*Services.ProdListResponse)
	result.Data=models
}


func (this *ProdsWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	cmdName:=req.Service()+" . "+req.Endpoint()
	configA:=hystrix.CommandConfig{
		Timeout:               1000,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:   50,
		//RequestVolumeThreshold:   2,
		//SleepWindow:            5000,
	}
	hystrix.ConfigureCommand("getprods",configA)
	return hystrix.Do(cmdName, func() error {
		return this.Client.Call(ctx,req,rsp)
	}, func(e error) error {
		fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+"调用超时熔断接口"+"[wrapper] client request service: %s method: %s\n", req.Service(), req.Endpoint())
		defaultData(rsp)
		return nil
	})
}

func NewlogWrap(c client.Client) client.Client {
	return &ProdsWrapper{c}
}
