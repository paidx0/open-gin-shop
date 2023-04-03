package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"shop-srvs/inventory_srv/model"

	"shop-srvs/inventory_srv/global"
	"shop-srvs/inventory_srv/proto"
)

type InventoryServer struct{}

// SetInv 设置库存，更新库存
func (*InventoryServer) SetInv(_ context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	if result := global.DB.Save(&inv); result.Error != nil {
		return nil, result.Error
	}

	return &emptypb.Empty{}, nil
}

// InvDetail 商品库存详情
func (*InventoryServer) InvDetail(_ context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "未找到商品库存详情")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// Sell 提交订单时，库存预先扣减相应数量，Redis分布式锁
func (*InventoryServer) Sell(_ context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	// 开始事务
	tx := global.DB.Begin()

	var details []model.GoodsDetail
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})

		// 获取Redis互斥锁
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		// 释放 Redis锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}

	// 保存库存订单扣减记录
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
		Detail:  details,
	}
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "保存库存扣减历史失败")
	}

	// 提交事务
	tx.Commit()

	return &emptypb.Empty{}, nil
}

// Reback 订单归还之前扣减的库存，订单超时、订单创建失败、手动取消订单
func (*InventoryServer) Reback(_ context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 开始事务
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	// 提交事务
	tx.Commit()
	return &emptypb.Empty{}, nil
}

// AutoReback Rocketmq 回查或收到回滚消息自动归还库存
func AutoReback(_ context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range msgs {
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		tx := global.DB.Begin()
		// 库存已扣减未归还的情况
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		for _, orderGood := range sellDetail.Detail {
			if result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGood.Num)); result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		if result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}

		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}

// // TrySell TCC 事务模式
// func (*InventoryServer) TrySell(req *proto.SellInfo) (*emptypb.Empty, error) {
// 	client := goredislib.NewClient(&goredislib.Options{
// 		Addr: "127.0.0.1:6379",
// 	})
// 	pool := goredis.NewPool(client)
// 	rs := redsync.New(pool)
//
// 	tx := global.DB.Begin()
// 	for _, goodInfo := range req.GoodsInfo {
// 		var inv model.InventoryNew
// 		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
// 		if err := mutex.Lock(); err != nil {
// 			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
// 		}
//
// 		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
// 		}
// 		// 判断库存是否充足
// 		if inv.Stocks < goodInfo.Num {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
// 		}
//
// 		inv.Freeze += goodInfo.Num
// 		tx.Save(&inv)
//
// 		if ok, err := mutex.Unlock(); !ok || err != nil {
// 			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
// 		}
// 	}
// 	tx.Commit()
// 	return &emptypb.Empty{}, nil
// }
//
// // ConfirmSell TCC 事务模式
// func (*InventoryServer) ConfirmSell(req *proto.SellInfo) (*emptypb.Empty, error) {
// 	client := goredislib.NewClient(&goredislib.Options{
// 		Addr: "127.0.0.1:6379",
// 	})
// 	pool := goredis.NewPool(client)
// 	rs := redsync.New(pool)
//
// 	tx := global.DB.Begin()
// 	for _, goodInfo := range req.GoodsInfo {
// 		var inv model.InventoryNew
// 		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
// 		if err := mutex.Lock(); err != nil {
// 			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
// 		}
//
// 		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
// 		}
//
// 		if inv.Stocks < goodInfo.Num {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
// 		}
//
// 		inv.Stocks -= goodInfo.Num
// 		inv.Freeze -= goodInfo.Num
// 		tx.Save(&inv)
//
// 		if ok, err := mutex.Unlock(); !ok || err != nil {
// 			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
// 		}
// 	}
// 	tx.Commit()
// 	return &emptypb.Empty{}, nil
// }
//
// // CancelSell TCC 事务模式
// func (*InventoryServer) CancelSell(req *proto.SellInfo) (*emptypb.Empty, error) {
// 	client := goredislib.NewClient(&goredislib.Options{
// 		Addr: "127.0.0.1:6379",
// 	})
// 	pool := goredis.NewPool(client)
// 	rs := redsync.New(pool)
//
// 	tx := global.DB.Begin()
// 	for _, goodInfo := range req.GoodsInfo {
// 		var inv model.InventoryNew
// 		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
// 		if err := mutex.Lock(); err != nil {
// 			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
// 		}
//
// 		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
// 		}
// 		if inv.Stocks < goodInfo.Num {
// 			tx.Rollback()
// 			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
// 		}
//
// 		inv.Freeze -= goodInfo.Num
// 		tx.Save(&inv)
//
// 		if ok, err := mutex.Unlock(); !ok || err != nil {
// 			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
// 		}
// 	}
// 	tx.Commit()
// 	return &emptypb.Empty{}, nil
// }
