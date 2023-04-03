package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Panicf("sentinel限流服务初始失败: %s", err.Error())
	}

	// 配置 sentinel限流规则
	rule := []*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              20,
			StatIntervalInMs:       6000,
		},
	}

	if _, err = flow.LoadRules(rule); err != nil {
		zap.S().Panicf("sentinel限流服务初始失败: %s", err.Error())
	}
}
