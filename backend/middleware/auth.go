package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Cookie中获取认证信息
		token, err := c.Cookie("auth_token")
		if err != nil || token != "admin_auth_token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
