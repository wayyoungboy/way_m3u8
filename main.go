package main

import (
	"github.com/gin-gonic/gin"
	"way_m3u8/controllers"
	"way_m3u8/worker"
)

func main() {
	worker.WorkInit()
	go worker.WorkRun()
	r := gin.Default()

	// 提供静态文件服务
	r.Static("/static", "./static") // 假设你的HTML文件在static目录下

	// ... 其他路由和中间件设置 ...
	r.POST("/url", controllers.StoreURL)

	r.Run(":2045") // 监听8080端口

}
