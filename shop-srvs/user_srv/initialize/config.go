package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop-srvs/user_srv/global"
)

// InitConfig 初始viper，并从nacos中获取配置信息
func InitConfig() {
	configFileName := fmt.Sprint("user_srv/config-debug.yaml")

	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("viper获取配置文件失败：%s", err.Error())
	}
	if err := v.Unmarshal(global.NacosConfig); err != nil {
		zap.S().Panicf("NacosConfig获取失败：%s", err.Error())
	}
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	// 开始从nacos拉取配置信息
	cc := constant.ClientConfig{
		NamespaceId:          global.NacosConfig.Namespace, // 命名空间
		TimeoutMs:            5000,                         // 请求超时
		CacheDir:             "tmp/nacos/cache",            // 缓存service信息的目录
		LogDir:               "tmp/nacos/log",              // 日志存储路径
		RotateTime:           "1h",                         // 日志轮转周期
		MaxAge:               3,                            // 日志最大文件数
		LogLevel:             "debug",                      // 日志默认级别
		NotLoadCacheAtStart:  true,                         // 启动时不读取缓存在CacheDir的service信息
		UpdateCacheWhenEmpty: true,                         // service返回列表空时不更新缓存，用于推空保护
	}

	// 至少一个ServerConfig
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	// 创建动态配置客户端
	// configClient, err := clients.CreateConfigClient(map[string]interface{}{
	// 	"serverConfigs": sc,
	// 	"clientConfig":  cc,
	// })
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		zap.S().Panicf("Nacos服务连接失败：%s", err.Error())
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Panicf("Nacos获取服务信息失败：%s", err.Error())
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("Nacos读取配置失败： %s", err.Error())
	}

	// fmt.Println(&global.ServerConfig)
}
