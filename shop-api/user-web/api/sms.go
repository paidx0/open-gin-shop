package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"shop-api/user-web/forms"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"shop-api/user-web/global"
)

// GenerateSmsCode 生成短信验证码
func GenerateSmsCode(witdh int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < witdh; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(10)])
	}
	return sb.String()
}

// SendSms 发送短信验证码
func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	if err != nil {
		zap.S().Errorf("[SendSms] sms连接失败：%s", err.Error())
	}
	// 6位短信验证码
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile
	request.QueryParams["SignName"] = "open-shop"                       // 阿里云验证过的项目名
	request.QueryParams["TemplateCode"] = "SMS_181850725"               // 阿里云的短信模板
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" // 短信模板中的验证码内容
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		zap.S().Errorf("[SendSms] 发送短信验证码失败：%s", err.Error())
	}

	// 验证码存 redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
