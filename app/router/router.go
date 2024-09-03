package router

import (
	"github.com/gantoho/go-img-sys/app/logic"
	"github.com/gantoho/go-img-sys/app/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()

	router.Use(middleware.Headers)

	router.GET("/f/:filename", logic.ReadFiles)

	v1 := router.Group("/v1")
	v1.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Go Yes!")
	})
	v1.GET("/all", logic.OpFiles)
	v1.GET("/get/:number", logic.BgimgNum)
	v1.GET("/bgimg", logic.Bgimg)
	v1.POST("/upload", logic.Upload)
	err := router.Run(":3128")
	if err != nil {
		return
	}
}
