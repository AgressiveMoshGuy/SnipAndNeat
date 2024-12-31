package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

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
		Name:        string(in.Name),
		Price:       in.Price,
	}
}

func (dto service) ToModel() *ozon_cli.ListTransactionsResultOperationService {
	return &ozon_cli.ListTransactionsResultOperationService{
		Name:  ozon_cli.TransactionOperationService(dto.Name),
		Price: dto.Price,
	}
}

func (db *DB) CreateService(ctx context.Context,
	in ozon_cli.ListTransactionsResultOperationService,
	operationID int64) (int64, error) {
	service := newService(in, operationID)

	if err := db.gdb.WithContext(ctx).
		Where("operation_id = ? AND name = ?", operationID, service.Name).
		Attrs(service).
		FirstOrCreate(&service).Error; err != nil {
		log.Error().Err(err)
		return 0, err
	}

	return service.ID, nil
}

func (db *DB) GetService(ctx context.Context, in ozon_cli.ListTransactionsResultOperationService) (*ozon_cli.ListTransactionsResultOperationService, error) {
	dto := &service{
		Name:  string(in.Name),
		Price: in.Price,
	}

	if err := db.gdb.WithContext(ctx).
		Where("name", dto.Name).
		First(dto).Error; err != nil {
		return nil, err
	}

	return dto.ToModel(), nil
}

func (db *DB) GetSumService(ctx context.Context, operationId int64) (float64, error) {
	dto := service{}

	if err := db.gdb.WithContext(ctx).Select("sum(price) as price, operation_id").
		Where("operation_id = ?", operationId).
		Group("operation_id").Find(&dto).Error; err != nil {
		return 0, fmt.Errorf("cannot get sum of services id=%d", operationId)
	}

	return dto.Price, nil
}

func (db *DB) GetSumServices(ctx context.Context, date time.Time) (map[string]float64, error) {
	dateFormatted := date.Format("2006-01-02 15:04:05")
	sumServices := make(map[string]float64)
	rows, err := db.gdb.WithContext(ctx).
		Table("services").
		Select("sn.description, SUM(services.price) AS price").
		Joins("left join operations o on services.operation_id = o.operation_id").
		Where("DATE(o.operation_date) = DATE(?)", dateFormatted).
		Joins("left join services_names sn on sn.name = services.name").
		Group("services.name").
		Rows()
	if err != nil {
		return nil, fmt.Errorf("cannot get sum of services for date=%v", date)
	}

	for rows.Next() {
		var (
			name  string
			price float64
		)
		if err := rows.Scan(&name, &price); err != nil {
			return nil, fmt.Errorf("cannot scan row: %v", err)
		}

		sumServices[name] = price
	}

	return sumServices, nil
}
