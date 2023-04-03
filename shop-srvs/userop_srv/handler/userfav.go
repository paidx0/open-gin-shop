package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/userop_srv/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/userop_srv/global"
	"shop-srvs/userop_srv/proto"
)

// GetFavList 收藏列表
func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var userFavs []model.UserFav
	result := global.DB.Where(&model.UserFav{User: req.UserId, Goods: req.GoodsId}).Find(&userFavs)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户未收藏商品")
	}

	var userFavList []*proto.UserFavResponse
	for _, userFav := range userFavs {
		userFavList = append(userFavList, &proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}

	return &proto.UserFavListResponse{
		Total: int32(result.RowsAffected),
		Data:  userFavList,
	}, nil
}

// AddUserFav 添加收藏
func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	userFav := &model.UserFav{
		User:  req.UserId,
		Goods: req.GoodsId,
	}

	if result := global.DB.Save(&userFav); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

// DeleteUserFav 删除收藏
func (*UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	if result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}

	return &emptypb.Empty{}, nil
}

// GetUserFavDetail 检查是否已经被收藏
func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Find(&model.UserFav{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}

	return &emptypb.Empty{}, nil
}
