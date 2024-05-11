package router

import (
	"github.com/gin-gonic/gin"
	"test/app/logic"
	"test/app/middleware"
)

func InitRouter() {
	router := gin.Default()

	router.Use(middleware.Headers)

	router.GET("/all", logic.OpFiles)

	router.GET("/get/:number", logic.BgimgNum)

	router.GET("/bgimg", logic.Bgimg)

	router.POST("/upload", logic.Upload)

	err := router.Run(":3128")
	if err != nil {
		return
	}
}
