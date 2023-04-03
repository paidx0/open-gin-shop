package middlewares

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"shop-api/userop-web/global"
	"shop-api/userop-web/models"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token过期")
	TokenNotValidYet = errors.New("token失效")
	TokenMalformed   = errors.New("token格式错误")
	TokenInvalid     = errors.New("token不存在")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.ServerConfig.JWTInfo.SigningKey),
	}
}

// JWTAuth 鉴权头部信息 x-token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			c.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, map[string]string{
					"msg": "授权已过期",
				})
			} else {
				c.JSON(http.StatusUnauthorized, "未登陆")
			}

			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// CreateToken 创建token
func (j *JWT) CreateToken(claims models.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析token
func (j *JWT) ParseToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}

}

// RefreshToken 刷新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
