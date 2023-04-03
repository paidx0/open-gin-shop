package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/oss-web/config"
)

var (
	Trans ut.Translator

	ServerConfig = &config.ServerConfig{}

	NacosConfig = &config.NacosConfig{}
)
