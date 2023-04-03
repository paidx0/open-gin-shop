package message

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/userop-web/api"
	"shop-api/userop-web/forms"
	"shop-api/userop-web/global"
	"shop-api/userop-web/models"
	"shop-api/userop-web/proto"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)

	// 非管理员无法查看所有
	request := &proto.MessageRequest{}
	if model.AuthorityId != 2 {
		request.UserId = int32(userId.(uint))
	}

	// 调用 grpc MessageList
	rsp, err := global.MessageClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		result = append(result, map[string]interface{}{
			"id":      value.Id,
			"user_id": value.UserId,
			"type":    value.MessageType,
			"subject": value.Subject,
			"message": value.Message,
			"file":    value.File,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  result,
	})
}

func New(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	messageForm := forms.MessageForm{}
	if err := ctx.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用 grpc CreateMessage
	rsp, err := global.MessageClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})
	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
