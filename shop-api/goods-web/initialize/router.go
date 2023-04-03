package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shop-api/goods-web/middlewares"
	"shop-api/goods-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 跨域 cors
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	return Router
}
