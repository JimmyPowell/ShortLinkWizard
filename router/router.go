package router

import (
	"shortlink/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8081"}, // 更新为前端实际地址
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	shortLinkController := controller.NewShortLinkController()

	// 创建短链接
	r.POST("/api/shorten", shortLinkController.CreateShortLink)

	// 访问短链接并重定向
	r.GET("/:code", shortLinkController.Redirect)

	return r
}
