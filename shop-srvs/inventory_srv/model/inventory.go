package model

type GoodsDetail struct {
	Goods int32
	Num   int32
}

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int comment '乐观锁'"`
}

type InventoryNew struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"`
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int comment '乐观锁'"`
	Freeze  int32 `gorm:"type:int comment '冻结库存,预扣减'"`
}

type Delivery struct {
	Goods   int32  `gorm:"type:int;index"`
	Nums    int32  `gorm:"type:int"`
	OrderSn string `gorm:"type:varchar(200)"`
	Status  string `gorm:"type:varchar(200) comment '1-等待支付 2-支付成功 3-失败'"`
}

type StockSellDetail struct {
	OrderSn string          `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status  int32           `gorm:"type:varchar(200) comment '1-已扣减 2-已归还'"`
	Detail  GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}
