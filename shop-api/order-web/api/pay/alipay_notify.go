package pay

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"shop-api/order-web/proto"

	"shop-api/order-web/global"
)

// Notify 支付完成后的回调通知
func Notify(ctx *gin.Context) {
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

	// 提取出支付宝返回的响应信息
	noti, err := client.GetTradeNotification(ctx.Request)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	// 调用 grpc UpdateOrderStatus
	_, err = global.OrderSrvClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: noti.OutTradeNo,
		Status:  string(noti.TradeStatus),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.String(http.StatusOK, "success")
}
