package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	"shop-api/goods-web/api"
	"shop-api/goods-web/forms"
	"shop-api/goods-web/global"
	"shop-api/goods-web/proto"
)

func List(ctx *gin.Context) {
	// 调用 grpc GetAllCategorysList
	r, err := global.GoodsSrvClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)

	if err = json.Unmarshal([]byte(r.JsonData), &data); err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	ctx.JSON(http.StatusOK, data)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用 grpc GetSubCategory
	r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	subCategorys := make([]interface{}, 0)
	for _, value := range r.SubCategorys {
		subCategorys = append(subCategorys, map[string]interface{}{
			"id":              value.Id,
			"name":            value.Name,
			"level":           value.Level,
			"parent_category": value.ParentCategory,
			"is_tab":          value.IsTab,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":              r.Info.Id,
		"name":            r.Info.Name,
		"level":           r.Info.Level,
		"is_tab":          r.Info.IsTab,
		"parent_category": r.Info.ParentCategory,
		"sub_categorys":   subCategorys,
	})
}

func New(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用 grpc CreateCategory
	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":     rsp.Id,
		"name":   rsp.Name,
		"parent": rsp.ParentCategory,
		"level":  rsp.Level,
		"is_tab": rsp.IsTab,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用 grpc DeleteCategory
	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Update(ctx *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	i, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}

	// 调用 grpc UpdateCategory
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
