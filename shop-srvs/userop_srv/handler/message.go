package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop-srvs/userop_srv/global"
	"shop-srvs/userop_srv/model"
	"shop-srvs/userop_srv/proto"
)

// MessageList 获取用户所有留言
func (*UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var messages []model.LeavingMessages
	result := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messages)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户未留言过")
	}

	var messageList []*proto.MessageResponse
	for _, message := range messages {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	return &proto.MessageListResponse{
		Total: int32(result.RowsAffected),
		Data:  messageList,
	}, nil
}

// CreateMessage 添加留言
func (*UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	message := &model.LeavingMessages{
		User:        req.UserId,
		MessageType: req.MessageType,
		Subject:     req.Subject,
		Message:     req.Message,
		File:        req.File,
	}

	if result := global.DB.Save(&message); result.Error != nil {
		return nil, result.Error
	}

	return &proto.MessageResponse{
		Id: message.ID,
	}, nil
}
