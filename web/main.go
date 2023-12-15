package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qx66/basicDiag/internal/service"
	"go.uber.org/zap"
)

func main() {
	//
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("初始化logger失败")
		return
	}
	
	//
	route := gin.New()
	svc, err := service.NewService(logger)
	if err != nil {
		fmt.Println("初始化服务失败, err: ", err)
		return
	}
	
	route.POST("/v1/hook/diag/web/report", svc.RecordReport)
	route.GET("/v1/hook/diag/web/report", svc.GetReport)
	
	err = route.Run(":20000")
	if err != nil {
		fmt.Println("启动服务失败, err: ", err)
	}
}
