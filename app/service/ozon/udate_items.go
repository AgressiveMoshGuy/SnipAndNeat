package ozon

import (
	"context"
	"strconv"

	"github.com/diphantxm/ozon-api-client/ozon"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func (z *OzonAPI) UpdateEANCodesWithItems(ctx context.Context) (*int64, error) {
	r, err := z.client.Products().GetListOfProducts(ctx, &ozon.GetListOfProductsParams{})
	if err != nil {
		return nil, err
	}

	var producIDs []int64
	lo.Map(r.Result.Items, func(v ozon.GetListOfProductsResultItem, _ int) []int64 {
		producIDs = append(producIDs, v.ProductId)
		return producIDs
	})
	result, err := z.client.Products().ListProductsByIDs(ctx, &ozon.ListProductsByIDsParams{
		ProductId: producIDs,
	})
	if err != nil {
		return nil, err
	}

	var reqMap = make(map[int64]int64, len(result.Items))
	for _, v := range result.Items {
		barcode, err := strconv.Atoi(v.Barcodes[0])
		if err != nil {
			continue
		}
		reqMap[v.SKU] = int64(barcode)
	}
	_, rows, err := z.db.WithTransaction(ctx, func(tx *gorm.DB, in any) (any, error) {
		return z.db.UpdateEANCodesBySKUs(tx, in.(map[int64]int64))
	}, reqMap)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}
