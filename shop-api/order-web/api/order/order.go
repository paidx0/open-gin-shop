package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"shop-api/order-web/forms"
	"shop-api/order-web/global"
	"shop-api/order-web/proto"

	"shop-api/order-web/api"
	"shop-api/order-web/models"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	request := proto.OrderFilterRequest{}

	// 如果是管理员用户则返回所有的订单
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	pagesInt, _ := strconv.Atoi(ctx.DefaultQuery("p", "0"))
	request.Pages = int32(pagesInt)

	perNumsInt, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))
	request.PagePerNums = int32(perNumsInt)

	// 调用 grpc OrderList
	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	orderList := make([]interface{}, 0)
	for _, item := range rsp.Data {
		orderList = append(orderList, map[string]interface{}{
			"id":       item.Id,
			"status":   item.Status,
			"pay_type": item.PayType,
			"user":     item.UserId,
			"post":     item.Post,
			"total":    item.Total,
			"address":  item.Address,
			"name":     item.Name,
			"mobile":   item.Mobile,
			"order_sn": item.OrderSn,
			"add_time": item.AddTime,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  orderList,
	})
}

func New(ctx *gin.Context) {
	orderForm := forms.CreateOrderForm{}
	if err := ctx.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidatorError(ctx, err)
	}

	// 调用 grpc CreateOrder
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(), "ginContext", ctx), &proto.OrderRequest{
		UserId:  int32(userId.(uint)),
		Name:    orderForm.Name,
		Mobile:  orderForm.Mobile,
		Address: orderForm.Address,
		Post:    orderForm.Post,
	})
	if err != nil {
		zap.S().Errorw("新建订单失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成支付宝链接
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 生成网页支付
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL              // 支付后回调链接
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL              // 返回链接
	p.Subject = "shop订单-" + rsp.OrderSn                                 // 标题
	p.OutTradeNo = rsp.OrderSn                                          // 订单号
	p.TotalAmount = strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64) // 金额
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"                            // 支付方式状态

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付链接失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"alipay_url": url.String(),
	})
}

func Detail(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	i, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "url格式出错",
		})
		return
	}
	request := proto.OrderRequest{
		Id: int32(i),
	}

	// 如果是管理员用户则返回所有的订单
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		request.UserId = int32(userId.(uint))
	}

	// 调用 grpc OrderDetail
	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &request)
	if err != nil {
		zap.S().Errorw("获取订单详情失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	goodsList := make([]interface{}, 0)
	for _, item := range rsp.Goods {
		goodsList = append(goodsList, map[string]interface{}{
			"id":    item.GoodsId,
			"name":  item.GoodsName,
			"image": item.GoodsImage,
			"price": item.GoodsPrice,
			"nums":  item.Nums,
		})
	}

	// 生成支付宝链接
	client, err := alipay.New(global.ServerConfig.AliPayInfo.AppID, global.ServerConfig.AliPayInfo.PrivateKey, false)
	if err != nil {
		zap.S().Errorw("实例化支付宝失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	err = client.LoadAliPayPublicKey(global.ServerConfig.AliPayInfo.AliPublicKey)
	if err != nil {
		zap.S().Errorw("加载支付宝的公钥失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 生成网页支付
	var p = alipay.TradePagePay{}
	p.NotifyURL = global.ServerConfig.AliPayInfo.NotifyURL                        // 支付后回调链接
	p.ReturnURL = global.ServerConfig.AliPayInfo.ReturnURL                        // 返回链接
	p.Subject = "shop订单-" + rsp.OrderInfo.OrderSn                                 // 标题
	p.OutTradeNo = rsp.OrderInfo.OrderSn                                          // 订单号
	p.TotalAmount = strconv.FormatFloat(float64(rsp.OrderInfo.Total), 'f', 2, 64) // 金额
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"                                      // 支付方式状态

	url, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Errorw("生成支付链接失败")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         rsp.OrderInfo.Id,
		"status":     rsp.OrderInfo.Status,
		"user":       rsp.OrderInfo.UserId,
		"post":       rsp.OrderInfo.Post,
		"total":      rsp.OrderInfo.Total,
		"address":    rsp.OrderInfo.Address,
		"name":       rsp.OrderInfo.Name,
		"mobile":     rsp.OrderInfo.Mobile,
		"pay_type":   rsp.OrderInfo.PayType,
		"order_sn":   rsp.OrderInfo.OrderSn,
		"goods":      goodsList,
		"alipay_url": url.String(),
	})
}
