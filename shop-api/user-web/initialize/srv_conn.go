package initialize

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"shop-api/user-web/global"
	"shop-api/user-web/proto"
)

// InitSrvConn 从Consul注册中心轮询算法获取服务rpc的连接信息，达到负载均衡
func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.UserSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}

	global.UserSrvClient = proto.NewUserClient(userConn)
}

// InitSrvConn2 单独的一个Web和一个rpc一对一的服务连接
func InitSrvConn2() {
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panicf("Consul服务连接失败：%s", err.Error())
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		zap.S().Panicf("Consul服务过滤失败：%s", err.Error())
	}

	userSrvHost := ""
	userSrvPort := 0

	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	// 连接grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】", "msg", err.Error())
	}

	global.UserSrvClient = proto.NewUserClient(userConn)
}
