package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

const (
	GetSumOperationsQuery = `
		select operation_id,
			item_sku,
			transaction_type,
			sum(amount) as         amount,
			sum(accruals_for_sale) accruals_for_sale,
			sum(sale_commission)   sale_commission
		from ozon_orders
			where operation_date >= ? and operation_date <= ?
		group by operation_id,transaction_type, item_sku
	`
)

type operation struct {
	OperationID          int64 `gorm:"column:operation_id;primary_key"`
	OperationType        string
	OperationDate        string
	OperationTypeName    string
	DeliveryCharge       float64
	ReturnDeliveryCharge float64
	AccrualsForSale      float64
	SaleCommission       float64
	Amount               float64
	TransactionType      string
	PostingID            int64         `gorm:"posting_id"`
	ItemSKU              int64         `gorm:"column:item_sku"`
	Services             pq.Int64Array `gorm:"services;type:integer[]"`
}

func (t operation) TableName() string {
	return "operations"
}

func newOperation(in ozon_cli.ListTransactionsResultOperation) operation {
	op := operation{
		OperationID:          in.OperationId,
		OperationType:        in.OperationType,
		OperationDate:        in.OperationDate,
		OperationTypeName:    in.OperationTypeName,
		DeliveryCharge:       in.DeliveryCharge,
		ReturnDeliveryCharge: in.ReturnDeliveryCharge,
		AccrualsForSale:      in.AccrualsForSale,
		SaleCommission:       in.SaleCommission,
		Amount:               in.Amount,
		TransactionType:      in.Type,
	}

	return op
}

// ToModel конвертация в models.User
func (dto operation) ToModel() ozon_cli.ListTransactionsResultOperation {
	out := ozon_cli.ListTransactionsResultOperation{
		OperationId:          dto.OperationID,
		OperationType:        dto.OperationType,
		OperationDate:        dto.OperationDate,
		OperationTypeName:    dto.OperationTypeName,
		DeliveryCharge:       dto.DeliveryCharge,
		ReturnDeliveryCharge: dto.ReturnDeliveryCharge,
		AccrualsForSale:      dto.AccrualsForSale,
		SaleCommission:       dto.SaleCommission,
		Amount:               dto.Amount,
		Items:                []ozon_cli.ListTransactionsResultOperationItem{{SKU: dto.ItemSKU}},
		Type:                 dto.TransactionType,
	}

	out.Services = make([]ozon_cli.ListTransactionsResultOperationService, len(dto.Services))
	for i, v := range dto.Services {
		out.Services[i].Name = ozon_cli.TransactionOperationService(fmt.Sprint(v))
	}

	return out
}

func (db *DB) CreateOperation(ctx context.Context, in ozon_cli.ListTransactionsResultOperation, postingID int64, itemSKU int64) error {
	op := newOperation(in)
	op.PostingID = postingID
	op.ItemSKU = itemSKU
	if err := db.gdb.WithContext(ctx).
		Where("operation_id = ?", op.OperationID).
		Attrs(operation{OperationID: op.OperationID}).
		FirstOrCreate(&op).Error; err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func (db *DB) UpdateOperation(ctx context.Context, operation ozon_cli.ListTransactionsResultOperation, services []int64) error {
	dto := newOperation(operation)
	dto.Services = services
	if err := db.gdb.WithContext(ctx).Where("operation_id = ?", dto.OperationID).
		Updates(&dto).Error; err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func (db *DB) GetOperationsFromTime(ctx context.Context, from, to time.Time) ([]ozon_cli.ListTransactionsResultOperation, error) {
	dto := make([]operation, 0)
	if err := db.gdb.WithContext(ctx).
		Where("operation_date > ? and operation_date < ?", from, to).
		Find(&dto).Error; err != nil {
		return nil, errors.Wrapf(err, "cannot get operations from %v to %v", from, to)
	}

	// services := make([]service, 0)
	// if err := db.gdb.WithContext(ctx).

	out := make([]ozon_cli.ListTransactionsResultOperation, len(dto))
	for i, v := range dto {
		out[i] = v.ToModel()
	}
	return out, nil
}

// запрос на получение списка сумм по транзакциям
func (db *DB) GetSumOperationsFromTime(ctx context.Context, from, to time.Time) ([]ozon_cli.ListTransactionsResultOperation, error) {
	dto := make([]operation, 0)
	if err := db.gdb.Raw(GetSumOperationsQuery, from, to).Scan(&dto).Error; err != nil {
		return nil, errors.Wrapf(err, "cannot get sum info from %v to %v", from, to)
	}
	out := make([]ozon_cli.ListTransactionsResultOperation, len(dto))
	for i, v := range dto {
		out[i] = v.ToModel()
	}
	return out, nil

}

func (db *DB) GetFirstOperation(ctx context.Context) (ozon_cli.ListTransactionsResultOperation, error) {
	dto := operation{}
	if err := db.gdb.Select("operation_date").Order("operation_date").First(&dto).Error; err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			return ozon_cli.ListTransactionsResultOperation{}, nil
		}
		return ozon_cli.ListTransactionsResultOperation{}, errors.Wrapf(err, "cannot get transaction")
	}
	return dto.ToModel(), nil
}

func (db *DB) GetLastOperation(ctx context.Context) (ozon_cli.ListTransactionsResultOperation, error) {
	dto := operation{}
	if err := db.gdb.Select("operation_date").Order("operation_date DESC").First(&dto).Error; err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			return ozon_cli.ListTransactionsResultOperation{}, nil
		}
		return ozon_cli.ListTransactionsResultOperation{}, errors.Wrapf(err, "cannot get transaction")
	}
	return dto.ToModel(), nil
}

func (db *DB) GetOperation(ctx context.Context, operationID int64) string {
	serviceName := struct {
		ID            int64  `gorm:"column:id;primary_key"`
		OperationName string `gorm:"column:name"`
		Description   string `gorm:"column:description"`
	}{}
	if err := db.gdb.Select("name").Table("services_names").Where("id = ?", operationID).First(&serviceName).Error; err != nil {
		if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
			return ""
		}
		return ""
	}
	return serviceName.OperationName
}
