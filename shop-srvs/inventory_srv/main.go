package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"shop-srvs/inventory_srv/global"
	"shop-srvs/inventory_srv/handler"
	"shop-srvs/inventory_srv/initialize"
	"shop-srvs/inventory_srv/proto"
	"shop-srvs/inventory_srv/utils"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")
	flag.Parse()

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)
	zap.S().Info("ip: ", *IP)
	if *IP != "0.0.0.0" {
		global.ServerConfig.Host = *IP
	}
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port: ", *Port)

	// 注册 grpc InventoryServer服务
	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port))
	if err != nil {
		zap.S().Panicf("grpc服务注册失败：%s", err.Error())
	}

	// 注册 健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册 Consul服务
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panicf("Consul服务连接失败：%s", err.Error())
	}
	// 对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	// 生成注册对象
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceID,
		Port:    *Port,
		Tags:    global.ServerConfig.Tags,
		Address: fmt.Sprintf("%s", global.ServerConfig.Host),
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Panicf("Consul服务注册失败：%s", err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			zap.S().Panicf("grpc服务监听失败：%s", err.Error())
		}
	}()

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"127.0.0.1:9876"}),
		consumer.WithGroupName("shop-inventory"),
	)
	// 接收到回滚消息
	if err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		zap.S().Errorf("rocketmq读取消息失败：%s", err.Error())
	}
	if err = c.Start(); err != nil {
		zap.S().Panicf("rocketmq启动失败：%s", err.Error())
	}

	// ctrl+c
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Consul用户服务下线
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("Consul用户服务下线失败")
	}
	zap.S().Info("Consul用户服务下线成功")
}
