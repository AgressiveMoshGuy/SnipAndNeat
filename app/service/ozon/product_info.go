package ozon

import (
	"context"
	"time"

	"github.com/diphantxm/ozon-api-client/ozon"
)

// получить список всех товаров
func (z *OzonAPI) VientoProductsInfo(ctx context.Context, in *ozon.GetListOfProductsParams) (*ozon.GetListOfProductsResponse, error) {
	// result, err := z.client.Products().GetListOfProducts(ctx, in)
	// if err != nil {
	// 	return nil, err
	// }

	t1 := "2024-07-25T10:43:06.51"
	t2 := "2024-08-25T10:43:06.51"
	t3, _ := time.Parse("2006-01-02T15:04:05", t1)
	t4, _ := time.Parse("2006-01-02T15:04:05", t2)

	result, err := z.client.Finance().GetTotalTransactionsSum(ctx, &ozon.GetTotalTransactionsSumParams{
		TransactionType: "all",
		Date: ozon.GetTotalTransactionsSumDate{
			From: t3,
			To:   t4,
		},
	})

	if err != nil {
		return nil, err
	}

	// for _, v := range result.Result.Items {
	// if err := z.db.CreateVientoProduct(ctx, models.VientoProduct{
	// 	ProductID: v.ProductId,
	// 	OfferID:   v.OfferId,
	// }); err != nil {
	// 	z.log.Err(err).Msgf("cannot create new Viento Product %d", v.ProductId)
	// 	return nil, fmt.Errorf("cannot create new Viento Product %w", err)
	// }
	z.log.Info().Msgf("result: %v", result)
	// }

	return nil, nil
}
