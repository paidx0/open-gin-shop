package handler

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/goods_srv/model"

	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/goods_srv/global"
	"shop-srvs/goods_srv/proto"
)

// GetAllCategorysList 获取所有的分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	var categorys []model.Category
	result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	if result.Error != nil {
		return nil, result.Error
	}

	jsonData, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{
		Total:    int32(result.RowsAffected),
		JsonData: string(jsonData),
	}, nil
}

// GetSubCategory 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	Info := &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}

	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse

	result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,
			Name:           subCategory.Name,
			Level:          subCategory.Level,
			IsTab:          subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryID,
		})
	}

	return &proto.SubCategoryListResponse{
		Total:        int32(result.RowsAffected),
		Info:         Info,
		SubCategorys: subCategoryResponse,
	}, nil
}

// CreateCategory 新建分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	cMap := map[string]interface{}{
		"name":   req.Name,
		"level":  req.Level,
		"is_tab": req.IsTab,
	}
	if req.Level != 1 {
		cMap["parent_category_id"] = req.ParentCategory
	}

	category := model.Category{}
	if result := global.DB.Model(category).Create(cMap); result.Error != nil {
		return nil, result.Error
	}

	return &proto.CategoryInfoResponse{
		Id:   category.ID,
		Name: category.Name,
	}, nil
}

// DeleteCategory 删除分类
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// UpdateCategory 修改分类
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}

	if result := global.DB.Save(&category); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
