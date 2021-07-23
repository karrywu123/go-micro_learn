package router

import (
	"github.com/gin-gonic/gin"
	"go-micro_learn/example3/protos/Services"
)

func NewGinRouter(prodService Services.ProdService) *gin.Engine {
	ginRouter:=gin.New()
	ginRouter.Use(IniMiddleware(prodService),ErrorMiddleware())
	v1group:=ginRouter.Group("/v1")
	{
		v1group.Handle("POST","/prods",GetProdsList)
		v1group.Handle("GET","/prods/:ProdId",GetProdsDetail)
	}
	return ginRouter
}
