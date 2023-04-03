package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shop-api/goods-web/models"
)

// IsAdminAuth 1-普通用户 2-管理员
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
