package initialize

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"shop-api/oss-web/middlewares"
	"shop-api/oss-web/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	// 测试 oss
	Router.LoadHTMLFiles(fmt.Sprintf("oss-web/templates/index.html"))
	Router.StaticFS("/static", http.Dir(fmt.Sprintf("oss-web/static")))
	Router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	// 跨域 cors
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/oss/v1")
	router.InitOssRouter(ApiGroup)

	return Router
}
