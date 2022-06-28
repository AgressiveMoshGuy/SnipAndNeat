package storage

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

type service struct {
	ID          int64 `gorm:"column:id;primary_key;auto_increment:true"`
	OperationID int64
	Name        string
	Price       float64
}

func (s service) TableName() string {
	return "services"
}

func newService(in ozon_cli.ListTransactionsResultOperationService, operationID int64) service {
	return service{
		OperationID: operationID,
		Name:        in.Name,
		Price:       in.Price,
	}
}

func (dto service) ToModel() *ozon_cli.ListTransactionsResultOperationService {
	return &ozon_cli.ListTransactionsResultOperationService{
		Name:  dto.Name,
		Price: dto.Price,
	}
}

func (db *DB) CreateService(ctx context.Context,
	in ozon_cli.ListTransactionsResultOperationService,
	operationID int64) (int64, error) {
	service := newService(in, operationID)
	if err := db.gorm.WithContext(ctx).Clauses(clause.OnConflict{
		OnConstraint: "oper_id_name",
		DoNothing:    true}).
		Save(&service).Error; err != nil {
		log.Error().Err(err)
		return 0, err
	}

	return service.ID, nil
}

func (db *DB) GetService(ctx context.Context, in ozon_cli.ListTransactionsResultOperationService) (*ozon_cli.ListTransactionsResultOperationService, error) {
	dto := &service{
		Name:  in.Name,
		Price: in.Price,
	}

	if err := db.gorm.WithContext(ctx).
		Where("name", dto.Name).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) GetSumService(ctx context.Context, operationId int64) (float64, error) {
	dto := service{}

	if err := db.gorm.WithContext(ctx).Select("sum(price) as price, operation_id").
		Where("operation_id = ?", operationId).
		Group("operation_id").Find(&dto).Error; err != nil {
		return 0, fmt.Errorf("cannot get sum of services id=%d", operationId)
	}

	return dto.Price, nil
}
