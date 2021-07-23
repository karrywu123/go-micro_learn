package router

import (
	"github.com/gin-gonic/gin"
	"go-micro_learn/example3/protos/Services"
	"fmt"
)

func IniMiddleware(prodService Services.ProdService)gin.HandlerFunc{
	return func(context *gin.Context) {
		context.Keys =make(map[string]interface{})
		context.Keys["prodservice"]=prodService //付值
		context.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc{
	return func(context *gin.Context) {
		defer func() {
			if r:=recover();r!=nil{
				context.JSON(500,gin.H{"status":fmt.Sprint("%s",r)})
			}
		}()
	}
}
