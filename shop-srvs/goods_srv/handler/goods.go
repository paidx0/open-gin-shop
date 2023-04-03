package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic/v7"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop-srvs/goods_srv/global"
	"shop-srvs/goods_srv/model"
	"shop-srvs/goods_srv/proto"
)

type GoodsServer struct{}

// GoodsList 商品列表
// 使用es搜索出商品实际的id，再到数据库中取出商品的信息返回
// 关键词筛选、新品筛选、热门筛选、价格区间筛选、商品分类筛选
func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// ES 和 MySQL 过滤规则
	q := elastic.NewBoolQuery()
	localDB := global.DB.Model(model.Goods{})
	if req.KeyWords != "" {
		// 多字段搜索
		q = q.Must(elastic.NewMultiMatchQuery(req.KeyWords, "name", "goods_brief"))
	}
	if req.IsHot {
		// localDB = global.DB.Where("is_hot = ?", req.IsHot)
		q = q.Filter(elastic.NewTermQuery("is_hot", req.IsHot))
	}
	if req.IsNew {
		// localDB = global.DB.Where("is_new = ?", req.IsNew)
		q = q.Filter(elastic.NewTermQuery("is_new", req.IsNew))
	}
	if req.Brand > 0 {
		// localDB = global.DB.Where("brands_id = ?", req.Brand)
		q = q.Filter(elastic.NewTermQuery("brands_id", req.Brand))
	}
	if req.PriceMin > 0 {
		// localDB = global.DB.Where("shop_price >= ?", req.PriceMin)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Gte(req.PriceMin))
	}
	if req.PriceMax > 0 {
		// localDB = global.DB.Where("shop_price <= ?", req.PriceMax)
		q = q.Filter(elastic.NewRangeQuery("shop_price").Lte(req.PriceMax))
	}

	// 通过category去过滤出子类商品列表
	categoryIds := make([]interface{}, 0)
	if req.TopCategory > 0 {
		var category model.Category
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		// 构造子分类查询语句
		var subQuery string
		switch {
		case category.Level == 1:
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		case category.Level == 2:
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		case category.Level == 3:
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}

		// 结果集
		type Result struct {
			ID int32
		}
		var results []Result
		if result := global.DB.Model(model.Category{}).Raw(subQuery).Scan(&results); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类为空")
		}

		// 生成terms查询
		for _, re := range results {
			categoryIds = append(categoryIds, re.ID)
		}
		q = q.Filter(elastic.NewTermsQuery("category_id", categoryIds...))
	}

	// 分页
	if req.Pages == 0 {
		req.Pages = 1
	}
	switch {
	case req.PagePerNums > 100:
		req.PagePerNums = 100
	case req.PagePerNums <= 0:
		req.PagePerNums = 10
	}
	// 执行 ES
	result, err := global.EsClient.Search().Index(model.EsGoods{}.GetIndexName()).Query(q).From(int(req.Pages)).Size(int(req.PagePerNums)).Do(context.Background())
	if err != nil {
		return nil, err
	}

	// 拿到实际的商品ID，再到数据库里查询信息返回出来
	goodsIds := make([]int32, 0)
	for _, value := range result.Hits.Hits {
		goods := model.EsGoods{}
		_ = json.Unmarshal(value.Source, &goods)
		goodsIds = append(goodsIds, goods.ID)
	}

	var goods []model.Goods
	re := localDB.Preload("Category").Preload("Brands").Find(&goods, goodsIds)
	if re.Error != nil {
		return nil, re.Error
	}

	var goodsinfo []*proto.GoodsInfoResponse
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsinfo = append(goodsinfo, &goodsInfoResponse)
	}

	return &proto.GoodsListResponse{
		Total: int32(result.Hits.TotalHits.Value),
		Data:  goodsinfo,
	}, nil
}

// GetGoodsDetail 商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods
	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	goodsInfoResponse := ModelToResponse(goods)

	return &goodsInfoResponse, nil
}

// CreateGoods 创建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	// 开启事务，保证和 ES的一致
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id:   goods.ID,
		Name: goods.Name,
	}, nil
}

// DeleteGoods 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	return &emptypb.Empty{}, nil
}

// UpdateGoods 修改商品
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	// 开启事务，保证和 ES的一致
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &emptypb.Empty{}, nil
}

// BatchGetGoods 下单等批量拉取商品信息
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	var goods []model.Goods
	result := global.DB.Where(req.Id).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	var goodInfo []*proto.GoodsInfoResponse
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodInfo = append(goodInfo, &goodsInfoResponse)
	}

	return &proto.GoodsListResponse{
		Total: int32(result.RowsAffected),
		Data:  goodInfo,
	}, nil
}
