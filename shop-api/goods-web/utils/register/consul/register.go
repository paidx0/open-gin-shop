package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

// Register 注册服务到Consul
func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Panicf("Consul服务连接失败：%s", err.Error())
	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Panicf("Consul服务注册失败：%s", err.Error())
	}
	return nil
}

// DeRegister 删除服务从Consul
func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorf("Consul服务连接失败：%s", err.Error())
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
