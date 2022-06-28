package storage

import (
	"context"

	"github.com/rs/zerolog/log"
)

type item struct {
	ID          int64 `gorm:"column:id;primary_key;auto_increment:true"`
	Name        string
	SKU         int64
	VientoID    int64 //`gorm:"column:viento_id"`
	Consumption int64
}

func (i item) TableName() string {
	return "items"
}

func newItem(in models.Item) item {
	return item{
		Name: in.Name,
		SKU:  in.SKU,
	}
}

func (dto item) ToModel() *models.Item {
	return &models.Item{
		ID:       dto.ID,
		Name:     dto.Name,
		SKU:      dto.SKU,
		VientoID: dto.VientoID,
	}
}

func (db *DB) CreateItem(ctx context.Context, in models.Item) (int64, error) {
	item := newItem(in)
	if err := db.gorm.WithContext(ctx).Where("sku=?", item.SKU).FirstOrCreate(&item).Error; err != nil {
		log.Error().Err(err)
		return 0, err
	}
	return item.ID, nil
}

func (db *DB) GetItemBySKU(ctx context.Context, in int64) (*models.Item, error) {
	dto := &item{
		SKU: in,
	}

	if err := db.gorm.WithContext(ctx).
		Where("sku = ?", dto.SKU).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) GetItem(ctx context.Context, id int64) (*models.Item, error) {
	dto := &item{}
	if err := db.gorm.WithContext(ctx).
		Where("id = ?", id).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) UpdateProductConsumpion(ctx context.Context, in models.Item, consumption int64) error {
	dto := newItem(in)
	if err := db.gorm.WithContext(ctx).Model(&dto).
		Where("sku = ?", in.SKU).
		UpdateColumn("consumption", consumption).Error; err != nil {
		return err
	}
	return nil
}
