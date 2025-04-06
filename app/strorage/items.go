package storage

import (
	"context"

	models "SnipAndNeat/generated"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type item struct {
	ID          int64 `gorm:"column:id;primary_key;auto_increment:true"`
	Name        string
	SKU         int64
	EAN         int64
	VientoID    int64 //`gorm:"column:viento_id"`
	Consumption int64
}

func (i item) TableName() string {
	return "items"
}

func newItem(in models.Item) item {
	return item{
		VientoID: in.VientoID.Value,
		Name:     in.Name.Value,
		SKU:      in.Sku.Value,
	}
}

func (dto item) ToModel() *models.Item {
	return &models.Item{
		ID:          models.NewOptInt64(dto.ID),
		Name:        models.NewOptString(dto.Name),
		Sku:         models.NewOptInt64(dto.SKU),
		Ean:         models.NewOptInt64(dto.EAN),
		VientoID:    models.NewOptInt64(dto.VientoID),
		Consumption: models.NewOptInt64(dto.Consumption),
	}
}

func (db *DB) CreateItem(ctx context.Context, in models.Item) (int64, error) {
	item := newItem(in)
	if err := db.gdb.WithContext(ctx).Where("sku=?", item.SKU).FirstOrCreate(&item).Error; err != nil {
		log.Error().Err(err)
		return 0, err
	}
	return item.ID, nil
}

func (db *DB) GetItemBySKU(ctx context.Context, in int64) (*models.Item, error) {
	dto := &item{
		SKU: in,
	}

	if err := db.gdb.WithContext(ctx).
		Where("sku = ?", dto.SKU).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) GetItem(ctx context.Context, id int64) (*models.Item, error) {
	dto := &item{}
	if err := db.gdb.WithContext(ctx).
		Where("id = ?", id).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) UpdateProductConsumpion(ctx context.Context, in models.Item, consumption int64) error {
	dto := newItem(in)
	if err := db.gdb.WithContext(ctx).Model(&dto).
		Where("sku = ?", dto.SKU).
		UpdateColumn("consumption", consumption).Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateEANCodesBySKUs(tx *gorm.DB, skusEANCodesMap map[int64]int64) (any, error) {
	for sku, ean := range skusEANCodesMap {
		if err := tx.Model(&item{}).
			Where("sku = ?", sku).
			UpdateColumn("ean", ean).Error; err != nil {
			return nil, err
		}
	}
	return nil, nil
}
