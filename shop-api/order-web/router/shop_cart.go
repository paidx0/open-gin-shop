package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/order-web/api/shop_cart"
	"shop-api/order-web/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", shop_cart.List)          // 购物车列表
		GoodsRouter.DELETE("/:id", shop_cart.Delete) // 删除购物车商品
		GoodsRouter.POST("", shop_cart.New)          // 添加购物车商品
		GoodsRouter.PATCH("/:id", shop_cart.Update)  // 修改购物车商品数量
	}
}
