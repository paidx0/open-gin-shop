package model

import "time"

type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index comment '用户'"`
	Goods   int32 `gorm:"type:int;index comment '商品'"`
	Nums    int32 `gorm:"type:int comment '件数'"`
	Checked bool  // 是否被选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel
	User    int32  `gorm:"type:int;index comment '用户'"`
	OrderSn string `gorm:"type:varchar(30);index comment '订单号'"`
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"`

	Status     string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo    string `gorm:"type:varchar(100) comment '交易号'"`
	OrderMount float32
	PayTime    *time.Time `gorm:"type:datetime comment '支付时间'"`

	Address      string `gorm:"type:varchar(100) comment '收件地址'"`
	SignerName   string `gorm:"type:varchar(20) comment '收件人'"`
	SingerMobile string `gorm:"type:varchar(11) comment '收件手机号'"`
	Post         string `gorm:"type:varchar(20) comment '备注留言'"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

type OrderGoods struct {
	BaseModel
	Order      int32  `gorm:"type:int;index comment '用户'"`
	Goods      int32  `gorm:"type:int;index comment '商品'"`
	GoodsName  string `gorm:"type:varchar(100);index comment '商品名'"`
	GoodsImage string `gorm:"type:varchar(200) comment '商品图片'"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int comment '件数'"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
