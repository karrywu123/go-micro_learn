package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProModel struct {
	ProdID int
	ProName string
}

type ProdsRequest struct {
	Size int `form:"size"`
}

func NewProd(id int,pname string) *ProModel {
	return &ProModel{ProdID:id,ProName:pname}
}

func NewProdList(n int) []*ProModel {
	ret :=make([]*ProModel,0)
	for i :=0;i<n;i++{
		ret=append(ret,NewProd(100+i,"prodname"+strconv.Itoa(100+i)))
	}
	return ret
}

func main() {
	//gin 路由
	ginRouter:=gin.New()
	v1group:=ginRouter.Group("/v1")
	{
		v1group.Handle("POST","/prods", func(context *gin.Context) {
			var pr  ProdsRequest
			err :=context.Bind(&pr)
			if err!=nil && pr.Size<=0{
				pr = ProdsRequest{Size:2}
			}
			context.JSON(200,gin.H{
				"data":NewProdList(pr.Size),
			})
		})

	}

	httpsrv := httpServer.NewServer(
		server.Name("prodservice"),
		server.Address(":8088"),
	)
	// Register consul
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs =[]string{"127.0.0.1:8500"}
	})
	hd := httpsrv.NewHandler(ginRouter)
	if err := httpsrv.Handle(hd); err != nil {
		logger.Fatal(err)
	}
	service := micro.NewService(
		micro.Server(httpsrv),
		micro.Registry(consulReg),
	)

	// Run service
	service.Init()

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
