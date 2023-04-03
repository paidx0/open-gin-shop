package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop-srvs/order_srv/global"
	"shop-srvs/order_srv/model"
	"shop-srvs/order_srv/proto"
)

type OrderServer struct{}

// CartItemList 获取购物车列表
func (*OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}

	var shopCartsInfo []*proto.ShopCartInfoResponse
	for _, shopCart := range shopCarts {
		shopCartsInfo = append(shopCartsInfo, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &proto.CartItemListResponse{
		Total: int32(result.RowsAffected),
		Data:  shopCartsInfo,
	}, nil
}

// CreateCartItem 添加购物车，没有就新建，有就合并
func (*OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart); result.RowsAffected == 1 {
		// 商品已经在购物车，则更新数量
		shopCart.Nums += req.Nums
	} else {
		// 不在就新添加
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}

	if result := global.DB.Save(&shopCart); result.Error != nil {
		return nil, result.Error
	}

	return &proto.ShopCartInfoResponse{
		Id: shopCart.ID,
	}, nil
}

// UpdateCartItem 修改购物车数量，更新是否选中
func (*OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.Checked = req.Checked

	if result := global.DB.Save(&shopCart); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

// DeleteCartItem 删除购物车
func (*OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	return &emptypb.Empty{}, nil
}

// OrderList 订单列表
func (*OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var orders []model.OrderInfo
	if result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders); result.Error != nil {
		return nil, result.Error
	}

	var OrderInfo []*proto.OrderInfoResponse
	for _, order := range orders {
		OrderInfo = append(OrderInfo, &proto.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)

	return &proto.OrderListResponse{
		Total: int32(total),
		Data:  OrderInfo,
	}, nil
}

// OrderDetail 订单详情
func (*OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	if result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&order); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	// 订单细节
	orderInfo := proto.OrderInfoResponse{
		Id:      order.ID,
		UserId:  order.User,
		OrderSn: order.OrderSn,

		PayType: order.PayType,
		Status:  order.Status,
		Total:   order.OrderMount,

		Address: order.Address,
		Name:    order.SignerName,
		Mobile:  order.SingerMobile,
		Post:    order.Post,
	}

	// 购买的商品列表
	var orderGoods []model.OrderGoods
	if result := global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods); result.Error != nil {
		return nil, result.Error
	}

	var OrderItem []*proto.OrderItemResponse
	for _, orderGood := range orderGoods {
		OrderItem = append(OrderItem, &proto.OrderItemResponse{
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsPrice: orderGood.GoodsPrice,
			GoodsImage: orderGood.GoodsImage,
			Nums:       orderGood.Nums,
		})
	}

	return &proto.OrderInfoDetailResponse{
		OrderInfo: &orderInfo,
		Goods:     OrderItem,
	}, nil
}

// CreateOrder 创建订单
func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	orderListener := OrderListener{Ctx: ctx}
	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithNameServer([]string{"127.0.0.1:9876"}),
	)
	if err != nil {
		zap.S().Errorf("rocketmq生成producer失败: %s", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("rocketmq启动producer失败: %s", err.Error())
		return nil, err
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}
	jsonString, _ := json.Marshal(order)

	// 回滚消息 half
	if _, err = p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", jsonString)); err != nil {
		zap.S().Errorf("rocketmq消息发送失败: %s", err.Error())
		return nil, status.Error(codes.Internal, "rocketmq消息发送失败")
	}

	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}

	return &proto.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderAmount,
	}, nil
}

// UpdateOrderStatus 修改订单状态
func (*OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}

// OrderTimeout 订单超时归还库存
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)

		fmt.Printf("获取到订单超时消息: %v\n", time.Now())
		// 查询订单的支付状态，已支付则跳过
		var order model.OrderInfo
		if result := global.DB.Model(model.OrderInfo{}).Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			// 修改订单的状态为已支付
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"127.0.0.1:9876"}))
			if err != nil {
				panic("生成producer失败")
			}
			if err = p.Start(); err != nil {
				panic("启动producer失败")
			}

			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			return consumer.ConsumeSuccess, nil
		}
	}
	return consumer.ConsumeSuccess, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	Ctx         context.Context
}

// ExecuteLocalTransaction rocketmq事务
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	parentSpan := opentracing.SpanFromContext(o.Ctx)

	var shopCarts []model.ShoppingCart
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts); result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	// 购物车选中商品数量
	var goodsIds []int32
	goodsNumsMap := make(map[int32]int32)
	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods", opentracing.ChildOf(parentSpan.Context()))
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()

	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})

		goodsInvInfo = append(goodsInvInfo, &proto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 预减库存
	queryInvSpan := opentracing.GlobalTracer().StartSpan("query_inv", opentracing.ChildOf(parentSpan.Context()))
	if _, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{OrderSn: orderInfo.OrderSn, GoodsInfo: goodsInvInfo}); err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	queryInvSpan.Finish()

	// 创建订单
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	saveOrderSpan := opentracing.GlobalTracer().StartSpan("save_order", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Save(&orderInfo); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}
	saveOrderSpan.Finish()

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	// 添加到订单商品记录
	saveOrderGoodsSpan := opentracing.GlobalTracer().StartSpan("save_order_goods", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}
	saveOrderGoodsSpan.Finish()

	// 删除购物车选中的
	deleteShopCartSpan := opentracing.GlobalTracer().StartSpan("delete_shopcart", opentracing.ChildOf(parentSpan.Context()))
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}
	deleteShopCartSpan.Finish()

	// 发送延时消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"127.0.0.1:9876"}))
	if err != nil {
		panic("生成producer失败")
	}
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 30分钟订单超时
	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(16)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延时消息失败: %v\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	// 提交事务
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

// CheckLocalTransaction rocketmq回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	if result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo); result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}
