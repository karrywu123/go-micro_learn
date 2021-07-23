package protos

import (
	"context"
	"go-micro_learn/example3/protos/Services"
	"strconv"
)

type ProdService struct {

}

func NewProd(id int32,pname string) *Services.ProdModel {
	return &Services.ProdModel{ProdID:id,ProdName:pname}
}

func( *ProdService)GetProdsList(ctx context.Context,in *Services.ProdsRequest,res *Services.ProdListResponse)error{
	//time.Sleep(time.Second*1)
	models:=make([]*Services.ProdModel,0)
	var i int32
	for i=0;i<in.Size;i++{
		models=append(models,NewProd(100+i,"prodname"+strconv.Itoa(100+int(i))))
	}
	res.Data=models
	return nil
}

func (*ProdService)GetProdsDetail(ctx context.Context, req  *Services.ProdsRequest,rsp  *Services.ProdDetailResponse) error{
	//time.Sleep(time.Second*1)
	rsp.Data=NewProd(req.ProdId,"测试商品"+strconv.Itoa(int(req.ProdId)))
	return nil
}
