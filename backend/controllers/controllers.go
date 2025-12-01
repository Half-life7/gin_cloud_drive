package controllers

import (
	"fmt"
	"gin_cloud_drive/backend/config"
	"gin_cloud_drive/backend/utils"
	"gin_cloud_drive/logger"
	"gin_cloud_drive/system"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Home 首页
func Home(c *gin.Context) {
	c.File("./frontend/index.html")
}

// About 关于页面
func About(c *gin.Context) {
	c.File("./frontend/about.html")
}

// SystemStatus 系统状态页面
func SystemStatus(c *gin.Context) {
	c.File("./frontend/system.html")
}

// Login 登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(ip, userAgent, "登录失败", fmt.Sprintf("请求参数错误: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	cfg := config.GetConfig()
	if req.Username == cfg.User.AdminUsername && req.Password == cfg.User.AdminPassword {
		// 设置认证Cookie
		c.SetCookie("auth_token", "admin_auth_token", 3600, "/", "", false, false)
		logger.LogUserOperation(ip, userAgent, "登录成功", fmt.Sprintf("用户 %s 登录成功", req.Username))
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
		})
	} else {
		logger.LogError(ip, userAgent, "登录失败", fmt.Sprintf("用户名或密码错误，尝试用户名: %s", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
	}
}

// Logout 退出登录
func Logout(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 删除认证Cookie
	c.SetCookie("auth_token", "", -1, "/", "", false, false)
	logger.LogUserOperation(ip, userAgent, "退出登录", "用户退出登录成功")
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "退出成功",
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	info := system.GetSystemInfo()
	c.JSON(http.StatusOK, info)
}

// GetSystemHistory 获取系统历史数据
func GetSystemHistory(c *gin.Context) {
	hours := 1
	if h := c.Query("hours"); h != "" {
		if _, err := fmt.Sscanf(h, "%d", &hours); err != nil {
			hours = 1
		}
	}

	// 限制时间范围
	if hours < 1 {
		hours = 1
	}
	if hours > 168 { // 7天
		hours = 168
	}

	data := system.GetHistoryData(hours)
	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"hours": hours,
	})
}

// ListFiles 列出文件
func ListFiles(c *gin.Context) {
	path := c.Query("path")
	sortBy := c.DefaultQuery("sort_by", "name")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	files, err := utils.ListFiles(path, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取文件列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": files,
	})
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.LogError(ip, userAgent, "上传文件失败", fmt.Sprintf("获取上传文件失败: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败",
		})
		return
	}
	defer file.Close()

	// 获取上传路径
	path := c.PostForm("path")
	uploadPath := config.GetConfig().File.UploadPath
	fullPath := filepath.Join(uploadPath, path, header.Filename)
	relativePath := filepath.Join(path, header.Filename)

	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		logger.LogError(ip, userAgent, "上传文件失败", fmt.Sprintf("创建目录失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败",
		})
		return
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.LogError(ip, userAgent, "上传文件失败", fmt.Sprintf("创建文件失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文件失败",
		})
		return
	}
	defer dst.Close()

	// 复制文件内容
	size, err := io.Copy(dst, file)
	if err != nil {
		logger.LogError(ip, userAgent, "上传文件失败", fmt.Sprintf("保存文件失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存文件失败",
		})
		return
	}

	logger.LogFileOperation(ip, userAgent, "上传文件", relativePath, size)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
	})
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 使用通配符参数，需要去掉前导斜杠
	filename := c.Param("filename")
	if filename != "" && filename[0] == '/' {
		filename = filename[1:]
	}
	// 将URL中的正斜杠转换为系统路径分隔符
	filename = filepath.FromSlash(filename)
	uploadPath := config.GetConfig().File.UploadPath

	// 获取文件大小
	fileInfo, err := os.Stat(filepath.Join(uploadPath, filename))
	var fileSize int64 = 0
	if err == nil {
		fileSize = fileInfo.Size()
	}

	logger.LogFileOperation(ip, userAgent, "下载文件", filename, fileSize)
	c.FileAttachment(filepath.Join(uploadPath, filename), filepath.Base(filename))
}

// PreviewFile 预览文件
func PreviewFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 使用通配符参数，需要去掉前导斜杠
	filename := c.Param("filename")
	if filename != "" && filename[0] == '/' {
		filename = filename[1:]
	}
	// 将URL中的正斜杠转换为系统路径分隔符
	filename = filepath.FromSlash(filename)
	uploadPath := config.GetConfig().File.UploadPath

	logger.LogFileOperation(ip, userAgent, "预览文件", filename, 0)
	c.File(filepath.Join(uploadPath, filename))
}

// RenameFile 重命名文件
func RenameFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	var req struct {
		OldPath string `json:"old_path"`
		NewName string `json:"new_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(ip, userAgent, "重命名文件失败", fmt.Sprintf("请求参数错误: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := utils.RenameFile(req.OldPath, req.NewName); err != nil {
		logger.LogError(ip, userAgent, "重命名文件失败", fmt.Sprintf("重命名文件失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "重命名文件失败",
		})
		return
	}

	logger.LogFileOperation(ip, userAgent, "重命名文件", fmt.Sprintf("%s -> %s", req.OldPath, req.NewName), 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件重命名成功",
	})
}

// MoveFile 移动文件
func MoveFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	var req struct {
		OldPath string `json:"old_path"`
		NewPath string `json:"new_path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(ip, userAgent, "移动文件失败", fmt.Sprintf("请求参数错误: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := utils.MoveFile(req.OldPath, req.NewPath); err != nil {
		logger.LogError(ip, userAgent, "移动文件失败", fmt.Sprintf("移动文件失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("移动文件失败: %v", err),
		})
		return
	}

	logger.LogFileOperation(ip, userAgent, "移动文件", fmt.Sprintf("%s -> %s", req.OldPath, req.NewPath), 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件移动成功",
	})
}

// DeleteFile 删除文件
func DeleteFile(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 使用通配符参数，需要去掉前导斜杠
	filename := c.Param("filename")
	if filename != "" && filename[0] == '/' {
		filename = filename[1:]
	}
	// 将URL中的正斜杠转换为系统路径分隔符
	filename = filepath.FromSlash(filename)

	if err := utils.DeleteFile(filename); err != nil {
		logger.LogError(ip, userAgent, "删除文件失败", fmt.Sprintf("删除文件失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("删除文件失败: %v", err),
		})
		return
	}

	logger.LogFileOperation(ip, userAgent, "删除文件", filename, 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件删除成功",
	})
}

// CreateDirectory 创建目录
func CreateDirectory(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	var req struct {
		Path string `json:"path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(ip, userAgent, "创建目录失败", fmt.Sprintf("请求参数错误: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := utils.CreateDirectory(req.Path); err != nil {
		logger.LogError(ip, userAgent, "创建目录失败", fmt.Sprintf("创建目录失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建目录失败",
		})
		return
	}

	logger.LogFileOperation(ip, userAgent, "创建目录", req.Path, 0)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "目录创建成功",
	})
}

// GetCarouselImages 获取轮播图图片
func GetCarouselImages(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	images, err := utils.ListFiles("carousel", "name", "asc")
	if err != nil {
		logger.LogError(ip, userAgent, "获取轮播图失败", fmt.Sprintf("获取轮播图失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取轮播图失败",
		})
		return
	}

	logger.LogAccess(ip, userAgent, "获取轮播图", fmt.Sprintf("获取到 %d 张轮播图", len(images)))
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": images,
	})
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	cfg := config.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"username": cfg.User.AdminUsername,
			"role":     "admin",
			"status":   "online",
		},
	})
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 这里可以添加更新用户信息的逻辑
	// 目前我们只支持单管理员，所以暂时不允许修改用户名
	if req.Username != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "用户名不允许修改",
		})
		return
	}

	// 如果提供了新密码，则更新密码
	if req.Password != "" {
		// 这里可以添加更新密码的逻辑
		// 目前我们只是返回成功，实际项目中应该更新配置文件或数据库
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "用户信息更新成功",
	})
}

// GetSettings 获取系统设置
func GetSettings(c *gin.Context) {
	cfg := config.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"server": gin.H{
				"port": cfg.Server.Port,
			},
			"file": gin.H{
				"max_size": cfg.File.MaxSize,
			},
			"system": gin.H{
				"interval": cfg.System.Interval,
			},
		},
	})
}

// UpdateSettings 更新系统设置
func UpdateSettings(c *gin.Context) {
	var req struct {
		Server struct {
			Port string `json:"port"`
		} `json:"server"`
		File struct {
			MaxSize int64 `json:"max_size"`
		} `json:"file"`
		System struct {
			Interval int `json:"interval"`
		} `json:"system"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 这里可以添加更新系统设置的逻辑
	// 目前我们只是返回成功，实际项目中应该更新配置文件或数据库

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "系统设置更新成功",
	})
}

// GetLogs 查询日志列表
func GetLogs(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	logger.LogAccess(ip, userAgent, "查询日志", "用户查询日志列表")

	// 解析查询参数
	var params logger.LogQueryParams

	// 解析页码和每页大小
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params.Page = page
	params.PageSize = pageSize

	// 解析时间范围
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02T15:04:05", startDateStr); err == nil {
			params.StartDate = startDate
		}
	}
	if endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02T15:04:05", endDateStr); err == nil {
			params.EndDate = endDate
		}
	}

	// 解析其他查询条件
	params.Level = c.Query("level")
	params.Type = c.Query("type")
	params.IP = c.Query("ip")
	params.Action = c.Query("action")
	params.File = c.Query("file")

	// 查询日志
	result, err := logger.QueryLogs(params)
	if err != nil {
		logger.LogError(ip, userAgent, "查询日志失败", fmt.Sprintf("查询日志失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("查询日志失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetLogStats 获取日志统计信息
func GetLogStats(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	logger.LogAccess(ip, userAgent, "获取日志统计", "用户获取日志统计信息")

	// 获取日志统计信息
	stats, err := logger.GetLogStats()
	if err != nil {
		logger.LogError(ip, userAgent, "获取日志统计失败", fmt.Sprintf("获取日志统计失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("获取日志统计失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// ClearOldLogs 清理旧日志
func ClearOldLogs(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()
	logger.LogAccess(ip, userAgent, "清理旧日志", "用户清理旧日志")

	// 解析清理天数
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	// 清理旧日志
	deletedCount, err := logger.ClearOldLogs(days)
	if err != nil {
		logger.LogError(ip, userAgent, "清理旧日志失败", fmt.Sprintf("清理旧日志失败: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": fmt.Sprintf("清理旧日志失败: %v", err),
		})
		return
	}

	logger.LogSystemOperation(ip, userAgent, "清理旧日志", fmt.Sprintf("成功清理了 %d 个旧日志文件", deletedCount))
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fmt.Sprintf("成功清理了 %d 个旧日志文件", deletedCount),
		"data": gin.H{
			"deleted_count": deletedCount,
		},
	})
}
