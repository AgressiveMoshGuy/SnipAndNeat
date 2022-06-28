package storage

import (
	models "SnipAndNeat/generated"
)

type vientoProduct struct {
	IsFBOVisible bool    `json:"is_fbo_visible"`
	Archived     bool    `json:"archived"`
	IsFBSVisible bool    `json:"is_fbs_visible"`
	IsDiscounted bool    `json:"is_discounted"`
	OfferID      string  `json:"offer_id"`
	ProductID    int64   `json:"product_id"`
	Price        float64 `json:"price"`
}

func (v vientoProduct) TableName() string {
	return "viento_products"
}

func (db *DB) NewVientoProduct(v vientoProduct) (*models.VientoProduct, error) {
	result := db.gdb.Create(&v)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

// аргегировать из таблицы общее потребление
// func (db *DB) Consumption(ctx context.Context)
