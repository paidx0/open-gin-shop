package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop-api/goods-web/global"
	"shop-api/goods-web/initialize"
	"shop-api/goods-web/utils"
	"shop-api/goods-web/utils/register/consul"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	Router := initialize.Routers()
	initialize.InitTrans("zh")
	initialize.InitSrvConn()
	initialize.InitSentinel()
	viper.AutomaticEnv()

	// 生产环境获取环境变量，获取空闲端口
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err != nil {
			zap.S().Panicf("端口获取失败：%s", err.Error())
		}
		global.ServerConfig.Port = port
	}

	// 注册 consul 服务
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	registerClient := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	err := registerClient.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panicf("Consul服务注册失败：%s", err.Error())
	}

	// 启动服务
	zap.S().Debugf("启动服务 端口： %d", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panicf("服务启动失败：%s", err.Error())
	}

	// ctrl+c
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Consul用户服务下线
	if err = registerClient.DeRegister(serviceId); err != nil {
		zap.S().Infof("Consul用户服务下线失败：%s", err.Error())
	}
	zap.S().Info("Consul用户服务下线成功")
}
