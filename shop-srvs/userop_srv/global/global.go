package global

import (
	"gorm.io/gorm"
	"shop-srvs/userop_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
