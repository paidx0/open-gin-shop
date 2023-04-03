package address

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/userop-web/api"
	"shop-api/userop-web/forms"
	"shop-api/userop-web/global"
	"shop-api/userop-web/models"
	"shop-api/userop-web/proto"
)

func List(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)

	request := &proto.AddressRequest{}
	// 非管理员无法查看所有
	if currentUser.AuthorityId != 2 {
		userId, _ := ctx.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	// 调用 grpc GetAddressList
	rsp, err := global.AddressClient.GetAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		result = append(result, map[string]interface{}{
			"id":            value.Id,
			"user_id":       value.UserId,
			"province":      value.Province,
			"city":          value.City,
			"district":      value.District,
			"address":       value.Address,
			"signer_name":   value.SignerName,
			"signer_mobile": value.SignerMobile,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  result,
	})
}

func New(ctx *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用 grpc CreateAddress
	userId, _ := ctx.Get("userId")
	rsp, err := global.AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(userId.(uint)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("新建地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用 grpc DeleteAddress
	_, err = global.AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	addressForm := forms.AddressForm{}
	if err = ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用 grpc UpdateAddress
	_, err = global.AddressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(i),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
