package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"go_Multi_learn/go-micro-prodservice/protos"
	"go_Multi_learn/go-micro-prodservice/protos/Services"
)



func main(){
	consulReg := consul.NewRegistry(registry.Addrs(":8500"))
	service:=micro.NewService(
		micro.Name("prodservice"),
		micro.Address(":8088"),
		micro.Registry(consulReg),
		)
	data:=Services.RegisterProdServiceHandler(service.Server(),new(protos.ProdService))
	fmt.Println(data)
	service.Init()
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
