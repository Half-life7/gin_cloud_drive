package routes

import (
	"gin_cloud_drive/backend/controllers"
	"gin_cloud_drive/backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 首页
	r.GET("/", controllers.Home)

	// API路由组
	api := r.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", controllers.Logout)
		}

		// 用户中心路由（需要认证）
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/info", controllers.GetUserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.GET("/settings", controllers.GetSettings)
			user.PUT("/settings", controllers.UpdateSettings)
		}

		// 文件管理路由
		file := api.Group("/file")
		{
			// 游客可访问的路由
			file.GET("/list", controllers.ListFiles)
			file.GET("/download/*filename", controllers.DownloadFile)
			file.GET("/preview/*filename", controllers.PreviewFile)
			file.GET("/carousel", controllers.GetCarouselImages)

			// 管理员可访问的路由（需要认证）
			adminFile := file.Group("/")
			adminFile.Use(middleware.AuthMiddleware())
			{
				adminFile.POST("/upload", controllers.UploadFile)
				adminFile.PUT("/rename", controllers.RenameFile)
				adminFile.PUT("/move", controllers.MoveFile)
				adminFile.DELETE("/delete/*filename", controllers.DeleteFile)
				adminFile.POST("/mkdir", controllers.CreateDirectory)
			}
		}

		// 系统状态路由
		system := api.Group("/system")
		{
			system.GET("/info", controllers.GetSystemInfo)
			system.GET("/history", controllers.GetSystemHistory)
		}

		// 日志管理路由（需要认证）
		log := api.Group("/log")
		log.Use(middleware.AuthMiddleware())
		{
			log.GET("/list", controllers.GetLogs)
			log.GET("/stats", controllers.GetLogStats)
			log.DELETE("/clear", controllers.ClearOldLogs)
		}
	}

	// 关于页面
	r.GET("/about", controllers.About)
	r.GET("/about/system", controllers.SystemStatus)
}
