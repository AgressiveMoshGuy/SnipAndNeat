package storage

import (
	"context"

	models "github.com/diphantxm/ozon-api-client/ozon"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

type posting struct {
	// gorm.Model
	ID             int64 `gorm:"column:id;primary_key;auto_increment:false"`
	DeliverySchema string
	OrderDate      string
	PostingNumber  string
	WarehouseID    int64
}

func (t posting) TableName() string {
	return "posts"
}

func newPosting(in models.ListTransactionsResultOperationPosting) posting {
	return posting{
		DeliverySchema: in.DeliverySchema,
		OrderDate:      in.OrderDate,
		PostingNumber:  in.PostingNumber,
		WarehouseID:    in.WarehouseId,
	}
}
func (dto posting) ToModel() *models.ListTransactionsResultOperationPosting {
	return &models.ListTransactionsResultOperationPosting{
		DeliverySchema: dto.DeliverySchema,
		OrderDate:      dto.OrderDate,
		PostingNumber:  dto.PostingNumber,
		WarehouseId:    dto.WarehouseID,
	}
}

func (db *DB) CreateOrGetPosting(ctx context.Context,
	in models.ListTransactionsResultOperationPosting) (
	int64, error) {
	posting := newPosting(in)

	if err := db.gorm.WithContext(ctx).
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoNothing: true}).
		Where("posting_number=?", in.PostingNumber).
		FirstOrCreate(&posting).
		Error; err != nil {
		log.Error().Err(err)
		return 0, err
	}
	return posting.ID, nil
}
