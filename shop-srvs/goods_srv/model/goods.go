package model

import (
	"context"
	"strconv"

	"gorm.io/gorm"
	"shop-srvs/goods_srv/global"
)

// Category 多级分类
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);comment '分类名';not null" json:"name"`
	Level            int32       `gorm:"type:int;comment '当前层级';default:1;not null" json:"level"`
	IsTab            bool        `gorm:"default:false;comment '是否展示';not null" json:"is_tab"`
	ParentCategoryID int32       `gorm:"comment '父级分类ID'" json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"-;foreignKey:ParentCategoryID;references:ID;comment '子级分类'" json:"sub_category"`
}

// Brands 品牌
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);comment '品牌名';not null"`
	Logo string `gorm:"type:varchar(200);comment '品牌logo';default:'';not null"`
}

// GoodsCategoryBrand 分类和品牌归类
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;comment '商品分类';index:idx_category_brand,unique"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;comment '绑定的品牌';index:idx_category_brand,unique"`
	Brands     Brands
}

func (GoodsCategoryBrand) TableName() string {
	return "goodcategorybrand"
}

// Banner 广告轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);comment '商品图片';not null"`
	Url   string `gorm:"type:varchar(200);comment '商品详情地址';not null"`
	Index int32  `gorm:"type:int;comment '轮播顺序';default:1;not null"`
}

// Goods 商品
type Goods struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;comment '商品分类';not null "`
	Category   Category
	BrandsID   int32 `gorm:"type:int;comment '绑定的品牌';not null "`
	Brands     Brands

	Name            string   `gorm:"type:varchar(50);comment '商品名';not null "`
	GoodsSn         string   `gorm:"type:varchar(50);comment '仓库中商品对应编号';not null "`
	ClickNum        int32    `gorm:"type:int;comment '查看数';default:0;not null "`
	SoldNum         int32    `gorm:"type:int;comment '销售数';default:0;not null "`
	FavNum          int32    `gorm:"type:int;comment '收藏数';default:0;not null "`
	MarketPrice     float32  `gorm:"type:int;comment '日常价';not null "`
	ShopPrice       float32  `gorm:"type:int;comment '优惠价';not null "`
	GoodsBrief      string   `gorm:"type:varchar(100);comment '商品描述';not null "`
	Images          GormList `gorm:"type:varchar(1000);comment '商品展示轮播图片';not null "`
	DescImages      GormList `gorm:"type:varchar(1000);comment '商品详情展示图片';not null "`
	GoodsFrontImage string   `gorm:"type:varchar(200);comment '商品图片';not null "`

	OnSale   bool `gorm:"default:false;comment '是否上架';not null "`
	ShipFree bool `gorm:"default:false;comment '是否打折';not null "`
	IsNew    bool `gorm:"default:false;comment '是否新品';not null "`
	IsHot    bool `gorm:"default:false;comment '是否热销';not null "`
}

// 事务钩子函数，保证 MySQL和 ES一致性

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	if _, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background()); err != nil {
		return err
	}

	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	if _, err = global.EsClient.Update().Index(esModel.GetIndexName()).Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background()); err != nil {
		return err
	}

	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	if _, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background()); err != nil {
		return err
	}

	return nil
}
