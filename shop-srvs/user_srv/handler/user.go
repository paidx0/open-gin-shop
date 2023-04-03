package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/user_srv/global"
	"shop-srvs/user_srv/model"
	"shop-srvs/user_srv/proto"
)

type UserServer struct{}

// GetUserList 分页获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	// 分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}
	return rsp, nil
}

// GetUserByMobile 通过手机号码查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// GetUserById 通过id查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// CreateUser 新建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	// 密码加密
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToResponse(user)
	return &userInfoRsp, nil
}

// UpdateUser 更新用户 昵称、生日、性别
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

// CheckPassWord 密码验证
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
