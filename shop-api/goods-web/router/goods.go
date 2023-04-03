package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-web/middlewares"

	"shop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	{
		GoodsRouter.GET("", goods.List)
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)
		GoodsRouter.GET("/:id", goods.Detail)
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)
		GoodsRouter.GET("/:id/stocks", goods.Stocks)

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus)
	}
}
