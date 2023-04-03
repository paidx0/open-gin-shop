package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/goods_srv/global"
	"shop-srvs/goods_srv/model"
	"shop-srvs/goods_srv/proto"
)

// CategoryBrandList 分类品牌的绑定关系列表
func (s *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total)

	var categoryBrands []model.GoodsCategoryBrand
	if result := global.DB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrands); result.Error != nil {
		return nil, result.Error
	}

	var categoryResponses []*proto.CategoryBrandResponse
	for _, categoryBrand := range categoryBrands {
		categoryResponses = append(categoryResponses, &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.Category.ID,
				Name:           categoryBrand.Category.Name,
				Level:          categoryBrand.Category.Level,
				IsTab:          categoryBrand.Category.IsTab,
				ParentCategory: categoryBrand.Category.ParentCategoryID,
			},
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.Brands.ID,
				Name: categoryBrand.Brands.Name,
				Logo: categoryBrand.Brands.Logo,
			},
		})
	}

	return &proto.CategoryBrandListResponse{
		Total: int32(total),
		Data:  categoryResponses,
	}, nil
}

// GetCategoryBrandList 根据品牌分类获取轮播图列表
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	var category model.Category
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var categoryBrands []model.GoodsCategoryBrand
	result := global.DB.Preload("Brands").Where(&model.GoodsCategoryBrand{CategoryID: req.Id}).Find(&categoryBrands)
	if result.RowsAffected <= 0 {
		return nil, status.Error(codes.NotFound, "商品分类未找到相关品牌")
	}

	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brands.ID,
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		})
	}

	return &proto.BrandListResponse{
		Total: int32(result.RowsAffected),
		Data:  brandInfoResponses,
	}, nil
}

// CreateCategoryBrand 添加分类品牌的绑定及轮播图
func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId,
		BrandsID:   req.BrandId,
	}

	if result := global.DB.Save(&categoryBrand); result.Error != nil {
		return nil, result.Error
	}
	return &proto.CategoryBrandResponse{
		Id: categoryBrand.ID,
	}, nil
}

// DeleteCategoryBrand 删除分类品牌的绑定及轮播图
func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	return &emptypb.Empty{}, nil
}

// UpdateCategoryBrand 修改分类品牌的绑定及轮播图
func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.GoodsCategoryBrand
	if result := global.DB.First(&categoryBrand, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	categoryBrand.CategoryID = req.CategoryId
	categoryBrand.BrandsID = req.BrandId

	if result := global.DB.Save(&categoryBrand); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
