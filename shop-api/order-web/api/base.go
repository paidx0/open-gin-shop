package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop-api/order-web/global"
)

// RemoveTopStruct 去掉结构体字段前缀
func RemoveTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// HandleGrpcErrorToHttp 将grpc的errCode转换成http statusCode
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg:": "内部错误",
			})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
		case codes.Unavailable:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "用户服务不可用",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": e.Code(),
			})
		}
		return
	}
}

// HandleValidatorError Validator验证不合法
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}
