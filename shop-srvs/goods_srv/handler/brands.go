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

// BrandList 品牌列表和logo
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brands []model.Brands
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}

	return &proto.BrandListResponse{
		Total: int32(total),
		Data:  brandResponses,
	}, nil
}

// CreateBrand 新建品牌
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	if result := global.DB.Save(brand); result.Error != nil {
		return nil, result.Error
	}

	return &proto.BrandInfoResponse{
		Id: brand.ID,
	}, nil
}

// DeleteBrand 删除品牌
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateBrand 修改品牌
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	if result := global.DB.Save(&brands); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
