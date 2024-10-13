package storage

import (
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

// аргегировать из таблицы общее потребление
// func (db *DB) Consumption(ctx context.Context)
