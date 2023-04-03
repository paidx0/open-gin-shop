package global

import (
	"gorm.io/gorm"
	"shop-srvs/order_srv/config"
	"shop-srvs/order_srv/proto"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig

	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
)
