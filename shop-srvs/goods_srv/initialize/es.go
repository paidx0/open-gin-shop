package initialize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"shop-srvs/goods_srv/global"
	"shop-srvs/goods_srv/model"
)

func InitEs() {
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "shop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		zap.S().Panicf("elastic服务初始失败：%s", err.Error())
	}

	// 检查索引是否存在，不存在新建
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		zap.S().Panicf("elastic服务初始失败：%s", err.Error())
	}
	if !exists {
		_, err = global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			zap.S().Panicf("elastic新建索引失败：%s", err.Error())
		}
	}
}
