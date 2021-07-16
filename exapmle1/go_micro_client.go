package main

import (
	"context"
	"fmt"
	myhttp "github.com/asim/go-micro/plugins/client/http/v3"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/selector"
	"io/ioutil"
	"net/http"
)
//go自带http调用
func callApi(addr,path, method string)(string ,error){
	req,_:=http.NewRequest(method,"http://"+addr+path,nil)
	client:=http.DefaultClient
	res,err:=client.Do(req)
	if err!=nil{
		return "",err
	}
	defer res.Body.Close()
	buf,_:=ioutil.ReadAll(res.Body)
	return string(buf),nil
}


//原生包调用
func callApi2(s selector.Selector){

	myclient:=myhttp.NewClient(
		       client.Selector(s),
		       client.ContentType("application/json"),

		       )

	req:=myclient.NewRequest("prodservice","/v1/prods", map[string]interface{}{"size":4})

	var  rsp map[string]interface{}

	err:=myclient.Call(context.Background(),req,&rsp)
	if err != nil{
		logger.Fatal(err)
	}
	fmt.Println(rsp["data"])
}


func main() {
	// Register consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs =[]string{"127.0.0.1:8500"}
	})

	reg2:=consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
	getService,err:=reg.GetService("prodservice")
	if err!=nil{
		logger.Fatal(err)
	}
	next:=selector.Random(getService)
	node,err:=next()
	if err!=nil{
		logger.Fatal(err)
	}

	callRes,err :=callApi(node.Address,"/v1/prods","POST")
	if err!=nil{
		logger.Fatal(err)
	}

	fmt.Println("======自建请求=====")
	fmt.Println(callRes)
	fmt.Println("====原生包请求===============")
	mySelector :=selector.NewSelector(

		selector.Registry(reg2),
		selector.SetStrategy(selector.Random),
		)
	callApi2(mySelector)
}
