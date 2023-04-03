package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"shop-api/user-web/middlewares"
	"shop-api/user-web/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop-api/user-web/forms"
	"shop-api/user-web/global"
	"shop-api/user-web/global/reponse"
	"shop-api/user-web/proto"
)

// GetUserList 用户列表
func GetUserList(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	pnInt, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	pSizeInt, _ := strconv.Atoi(ctx.DefaultQuery("psize", "10"))

	// 调用grpc GetUserList
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := reponse.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			// Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: reponse.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		result = append(result, user)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  result,
	})
}

// PassWordLogin 登录验证
func PassWordLogin(c *gin.Context) {
	passwordLoginForm := forms.PassWordLoginForm{}
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 检查验证码
	if store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 调用grpc GetUserByMobile 先查找用户
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	}

	// 调用 grpc CheckPassWord 密码验证
	passRsp, pasErr := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.PassWord,
		EncryptedPassword: rsp.PassWord,
	})
	if pasErr != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"password": "登录失败",
		})
	}
	if passRsp.Success {
		// 生成token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint(rsp.Id),
			NickName:    rsp.NickName,
			AuthorityId: uint(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*60*24*30,
				Issuer:    "paidx0",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         rsp.Id,
			"nick_name":  rsp.NickName,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	} else {
		c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "登录失败",
		})
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 从 Redis 取出验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码不存在",
		})
		return
	}
	if value != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	// 调用 grpc CreateUser
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24*30,
			Issuer:    "paidx0",
		},
	}

	// 生成token
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

// GetUserDetail 用户信息
func GetUserDetail(ctx *gin.Context) {
	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	// 调用 grpc GetUserById
	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":     rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}

// UpdateUser 更新用户 昵称、生日、性别
func UpdateUser(ctx *gin.Context) {
	updateUserForm := forms.UpdateUserForm{}
	if err := ctx.ShouldBind(&updateUserForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)

	loc, _ := time.LoadLocation("Local")
	birthDay, _ := time.ParseInLocation("2006-01-02", updateUserForm.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: updateUserForm.Name,
		Gender:   updateUserForm.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
