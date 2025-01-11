package router

import (
	"shortlink/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	shortLinkController := controller.NewShortLinkController()

	// 创建短链接
	r.POST("/api/shorten", shortLinkController.CreateShortLink)

	// 访问短链接并重定向
	r.GET("/:code", shortLinkController.Redirect)

	return r
}
