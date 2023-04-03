package handler

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// Paginate 分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GenerateOrderSn 生成订单号，年月日时分秒+用户id+2位随机数
func GenerateOrderSn(userId int32) string {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(
		"%d%d%d%d%d%d%d%d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(),
		userId, rand.Intn(90)+10,
	)
}
