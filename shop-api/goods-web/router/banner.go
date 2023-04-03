package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-web/api/banners"
	"shop-api/goods-web/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners").Use(middlewares.Trace())
	{
		BannerRouter.GET("", banners.List)
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete)
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update)
	}
}
