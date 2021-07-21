package protos

import (
	"context"
	"protos/Services"
	"strconv"
)

type ProdService struct {

}

func NewProd(id int32,pname string) *Services.ProdModel {
	return &Services.ProdModel{ProdID:id,ProdName:pname}
}

func( *ProdService)GetProdsList(ctx context.Context,in *Services.ProdsRequest,res *Services.ProdListResponse)error{
	models:=make([]*Services.ProdModel,0)
	var i int32
	for i=0;i<in.Size;i++{
		models=append(models,NewProd(100+i,"prodname"+strconv.Itoa(100+int(i))))
	}
	res.Data=models
	return nil
}
