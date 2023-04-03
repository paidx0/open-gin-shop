package shop_cart

import (
	"context"
	"net/http"
	"strconv"

	"go.uber.org/zap"
	"shop-api/order-web/api"
	"shop-api/order-web/forms"
	"shop-api/order-web/proto"

	"github.com/gin-gonic/gin"
	"shop-api/order-web/global"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	// 调用 grpc CartItemList
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{Id: int32(userId.(uint))})
	if err != nil {
		zap.S().Errorw("[List] 查询 【购物车列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ids := make([]int32, 0)
	for _, item := range rsp.Data {
		ids = append(ids, item.GoodsId)
	}
	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 调用 grpc BatchGetGoods
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: ids})
	if err != nil {
		zap.S().Errorw("[List] 批量查询【商品列表】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		for _, good := range goodsRsp.Data {
			if good.Id == item.GoodsId {
				goodsList = append(goodsList, map[string]interface{}{
					"id":         item.Id,
					"goods_id":   item.GoodsId,
					"good_name":  good.Name,
					"good_image": good.GoodsFrontImage,
					"good_price": good.ShopPrice,
					"nums":       item.Nums,
					"checked":    item.Checked,
				})
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  goodsList,
	})
}

func New(ctx *gin.Context) {
	itemForm := forms.ShopCartItemForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用 grpc GetGoodsDetail
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: itemForm.GoodsId})
	if err != nil {
		zap.S().Errorw("[List] 查询【商品信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 调用 grpc InvDetail
	invRsp, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: itemForm.GoodsId,
	})
	if err != nil {
		zap.S().Errorw("[List] 查询【库存信息】失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if invRsp.Num < itemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	// 调用 grpc CreateCartItem
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		GoodsId: itemForm.GoodsId,
		UserId:  int32(userId.(uint)),
		Nums:    itemForm.Nums,
	})
	if err != nil {
		zap.S().Errorw("添加到购物车失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(ctx *gin.Context) {
	i, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	itemForm := forms.ShopCartItemUpdateForm{}
	if err := ctx.ShouldBindJSON(&itemForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
		Nums:    itemForm.Nums,
		Checked: false,
	}
	if itemForm.Checked != nil {
		request.Checked = *itemForm.Checked
	}

	// 调用 grpc UpdateCartItem
	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("更新购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	i, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	// 调用 grpc DeleteCartItem
	userId, _ := ctx.Get("userId")
	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(i),
	})
	if err != nil {
		zap.S().Errorw("删除购物车记录失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
