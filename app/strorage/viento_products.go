package storage

import (
	"context"
	models "xls_importer/models/ozon_models"
)

type vientoProduct struct {
	IsFBOVisible bool
	Archived     bool
	IsFBSVisible bool
	IsDiscounted bool
	OfferID      string
	ProductID    int64
	Price        float64
}

func (v vientoProduct) TableName() string {
	return "viento_products"
}

func newVientoProduct(in models.VientoProduct) vientoProduct {
	return vientoProduct{
		ProductID:    in.ProductID,
		OfferID:      in.OfferID,
		IsFBOVisible: in.IsFBOVisible,
		IsFBSVisible: in.IsFBSVisible,
		Archived:     in.Archived,
		IsDiscounted: in.IsDiscounted,
		Price:        in.Price,
	}
}

func (dto vientoProduct) ToModel() *models.VientoProduct {
	return &models.VientoProduct{
		ProductID:    dto.ProductID,
		OfferID:      dto.OfferID,
		IsFBOVisible: dto.IsFBOVisible,
		IsFBSVisible: dto.IsFBSVisible,
		Archived:     dto.Archived,
		IsDiscounted: dto.IsDiscounted,
		Price:        dto.Price,
	}
}

func (db *DB) CreateVientoProduct(ctx context.Context, in models.VientoProduct) error {
	dto := newVientoProduct(in)
	if err := db.gorm.WithContext(ctx).FirstOrCreate(&dto).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) GetVientoProduct(ctx context.Context, productID int64) (*models.VientoProduct, error) {
	var dto vientoProduct
	if err := db.gorm.WithContext(ctx).Where("product_id = ?", productID).Find(&dto).Error; err != nil {
		return nil, err
	}
	return dto.ToModel(), nil
}

// аргегировать из таблицы общее потребление
// func (db *DB) Consumption(ctx context.Context)
