package main

import (
	"gin_cloud_drive/backend/config"
	"gin_cloud_drive/backend/routes"
	"gin_cloud_drive/logger"
	"gin_cloud_drive/system"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志系统
	if err := logger.InitLogger("./logs", "INFO"); err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}

	// 检测并创建carousel文件夹
	cfg := config.GetConfig()
	carouselPath := fmt.Sprintf("%s/carousel", cfg.File.UploadPath)
	if _, err := os.Stat(carouselPath); os.IsNotExist(err) {
		if err := os.MkdirAll(carouselPath, 0755); err != nil {
			logger.LogError("", "", "创建carousel文件夹失败", fmt.Sprintf("创建carousel文件夹失败: %v", err))
		} else {
			logger.LogSystemOperation("", "", "创建carousel文件夹", fmt.Sprintf("成功创建carousel文件夹: %s", carouselPath))
		}
	}

	// 初始化系统状态监控
	system.InitSystemMonitor()

	// 创建Gin引擎
	r := gin.Default()

	// 设置静态文件服务
	r.Static("/static", "./frontend")
	r.Static("/upload", "./upload")

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	port := config.GetConfig().Server.Port
	fmt.Printf("服务器运行在 http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
