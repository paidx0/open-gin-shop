package router

import (
	"github.com/gin-gonic/gin"
	"shop-api/goods-web/api/brands"
	"shop-api/goods-web/middlewares"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands").Use(middlewares.Trace())
	{
		BrandRouter.GET("", brands.BrandList)
		BrandRouter.DELETE("/:id", brands.DeleteBrand)
		BrandRouter.POST("", brands.NewBrand)
		BrandRouter.PUT("/:id", brands.UpdateBrand)
	}

	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand)
		CategoryBrandRouter.POST("", brands.NewCategoryBrand)
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)
	}
}
