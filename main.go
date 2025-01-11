package main

import (
	"log"
	"net/http"
	"shortlink/config"
	"shortlink/router"
)

func main() {
	// 初始化配置（包括数据库连接）
	config.InitConfig()

	// 初始化路由
	r := router.SetupRouter()

	// 启动服务
	port := "8081" // 使用 8081 端口
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
