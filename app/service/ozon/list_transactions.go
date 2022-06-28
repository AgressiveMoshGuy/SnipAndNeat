package ozon

import (
	"context"
	"fmt"
	"time"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

// список всех транзакций
func (z *OzonAPI) GetListTransaction(ctx context.Context, in *ozon_cli.ListTransactionsParams) (*ozon_cli.ListTransactionsResponse, error) {
	result, err := z.client.Finance().ListTransactions(ctx, in)
	if err != nil {
		z.log.Err(err).Msg("cannot get list transactions")
		return nil, err
	}

	for _, operation := range result.Result.Operations {
		postId, err := z.db.CreateOrGetPosting(ctx, operation.Posting)
		if err != nil {
			z.log.Err(err).Msg("cannot get or create posting")
			return nil, err
		}

		for _, item := range operation.Items {
			itemId, err := z.db.CreateItem(ctx, models.Item{
				Name: item.Name,
				SKU:  item.SKU,
			})
			if err != nil {
				z.log.Err(err).Msg("cannot create item")
				return nil, err
			}

			err = z.db.CreateOperation(ctx, operation, postId, itemId)
			if err != nil {
				z.log.Err(err).Msg("cannot create operation")
				return nil, err
			}
		}

		var services []int64
		for _, s := range operation.Services {
			id, err := z.db.CreateService(ctx, s, operation.OperationId)
			if err != nil {
				z.log.Err(err).Msg("cannot create service")
				return nil, err
			}
			services = append(services, id)
		}

		if err := z.db.UpdateOperation(ctx, operation, services); err != nil {
			z.log.Err(err).Msg("cannot update operation")
			return nil, err
		}
	}

	return result, nil
}

// полная история транзакций
func (z *OzonAPI) FullHistory(ctx context.Context, in *ozon_cli.ListTransactionsFilter) error {
	for !time.Now().Before(in.Date.To) {
		req := &ozon_cli.ListTransactionsParams{
			Filter:   *in,
			Page:     1,
			PageSize: 100,
		}
		list, err := z.GetListTransaction(ctx, req)
		if err != nil {
			z.log.Err(err).Msg("cannot create operation")
			return err
		}
		fmt.Println(len(list.Result.Operations))
		in.Date.To = in.Date.To.AddDate(0, 0, 1)
		in.Date.From = in.Date.From.AddDate(0, 0, 1)
	}

	return nil
}
