package model

const (
	LEAVING_MESSAGES = iota + 1
	COMPLAINT
	INQUIRY
	POST_SALE
	WANT_TO_BUY
)

type LeavingMessages struct {
	BaseModel
	User        int32  `gorm:"type:int;index"`
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100) comment '主题'"`

	Message string
	File    string `gorm:"type:varchar(200) comment 'oss上文件链接'"`
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

type Address struct {
	BaseModel
	User         int32  `gorm:"type:int;index comment '用户'"`
	Province     string `gorm:"type:varchar(10) comment '省'"`
	City         string `gorm:"type:varchar(10) comment '市'"`
	District     string `gorm:"type:varchar(20) comment '区'"`
	Address      string `gorm:"type:varchar(100) comment '地址'"`
	SignerName   string `gorm:"type:varchar(20) comment '收件人'"`
	SignerMobile string `gorm:"type:varchar(11) comment '手机号'"`
}

type UserFav struct {
	BaseModel
	User  int32 `gorm:"type:int;index:idx_user_goods,unique comment '用户'"`
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique comment '收藏商品'"`
}

func (UserFav) TableName() string {
	return "userfav"
}
