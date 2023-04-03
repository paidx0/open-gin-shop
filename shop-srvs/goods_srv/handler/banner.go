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

// BannerList 获取轮播图列表
func (s *GoodsServer) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	var banners []model.Banner
	result := global.DB.Find(&banners)
	if result.Error != nil {
		return nil, result.Error
	}

	var bannerResponses []*proto.BannerResponse
	for _, banner := range banners {
		bannerResponses = append(bannerResponses, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}

	return &proto.BannerListResponse{
		Total: int32(result.RowsAffected),
		Data:  bannerResponses,
	}, nil
}

// CreateBanner 添加轮播图
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{
		Image: req.Image,
		Url:   req.Url,
		Index: req.Index,
	}

	if result := global.DB.Save(&banner); result.Error != nil {
		return nil, result.Error
	}

	return &proto.BannerResponse{
		Id: banner.ID,
	}, nil
}

// DeleteBanner 删除轮播图
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

// UpdateBanner 修改轮播图
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner

	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	if result := global.DB.Save(&banner); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
