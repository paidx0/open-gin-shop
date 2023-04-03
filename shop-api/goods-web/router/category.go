package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-web/api/category"
	"shop-api/goods-web/middlewares"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("categorys").Use(middlewares.Trace())
	{
		CategoryRouter.GET("", category.List)
		CategoryRouter.DELETE("/:id", category.Delete)
		CategoryRouter.GET("/:id", category.Detail)
		CategoryRouter.POST("", category.New)
		CategoryRouter.PUT("/:id", category.Update)
	}
}
