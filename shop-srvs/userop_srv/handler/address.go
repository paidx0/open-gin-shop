package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"shop-srvs/userop_srv/global"
	"shop-srvs/userop_srv/model"
	"shop-srvs/userop_srv/proto"
)

// GetAddressList 查看地址
func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addresses)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户未添加地址")
	}

	var addressResponse []*proto.AddressResponse
	for _, address := range addresses {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}

	return &proto.AddressListResponse{
		Total: int32(result.RowsAffected),
		Data:  addressResponse,
	}, nil
}

// CreateAddress 新增地址
func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	address := &model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}

	if result := global.DB.Save(&address); result.Error != nil {
		return nil, result.Error
	}

	return &proto.AddressResponse{
		Id: address.ID,
	}, nil
}

// DeleteAddress 删除地址
func (*UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).Delete(&model.Address{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateAddress 修改地址
func (*UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	var address model.Address
	if result := global.DB.Where("id=? and user=?", req.Id, req.UserId).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if address.Province != "" {
		address.Province = req.Province
	}

	if address.City != "" {
		address.City = req.City
	}

	if address.District != "" {
		address.District = req.District
	}

	if address.Address != "" {
		address.Address = req.Address
	}

	if address.SignerName != "" {
		address.SignerName = req.SignerName
	}

	if address.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}

	if result := global.DB.Save(&address); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}
