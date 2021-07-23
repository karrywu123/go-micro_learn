package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-micro_learn/example3/protos/Services"
)


func PanicifError(err error){
	if err!=nil{
		panic(err)
	}
}

func GetProdsList(gincontext *gin.Context){
	prodService :=gincontext.Keys["prodservice"].(Services.ProdService)
	var prodReq Services.ProdsRequest
	err :=gincontext.Bind(&prodReq)
	if err!=nil{
		gincontext.JSON(500,gin.H{"status":err.Error()})

	}else {

		prodRes,_ := prodService.GetProdsList(context.Background(),&prodReq)
		gincontext.JSON(200,gin.H{"data":prodRes.Data})
		//加入熔断器超时出现异常
		//configA :=hystrix.CommandConfig{
		//	Timeout:               1000,
		//	MaxConcurrentRequests:  100,
		//	ErrorPercentThreshold:   25,
		//}
		////
		//hystrix.ConfigureCommand("getprods",configA)
		////
		//var prodRes *Services.ProdListResponse
		//err:=hystrix.Do("getprods", func() error {
		//	prodRes,err = prodService.GetProdsList(context.Background(),&prodReq)
		//	return err
		//}, func(e error) error {
		//	prodRes,err = defaultProds()
		//	return err
		//})
		////prodRes,_:=prodService.GetProdsList(context.Background(),&prodRes)
		//if err !=nil{
		//	gincontext.JSON(500,gin.H{"status":err.Error()})
		//}else {
		//	gincontext.JSON(200,gin.H{"data":prodRes.Data})
		//}
	}


}

func GetProdsDetail(ginctx *gin.Context){
	var proReq Services.ProdsRequest
	PanicifError(ginctx.BindUri(&proReq))
	prodService :=ginctx.Keys["prodservice"].(Services.ProdService)
	resp,_:=prodService.GetProdsDetail(context.Background(),&proReq)

	ginctx.JSON(200,gin.H{"data":&resp.Data})
}
