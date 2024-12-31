package ozon

import (
	"context"

	models "SnipAndNeat/generated"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

// список всех транзакций
func (z *OzonAPI) GetListTransaction(ctx context.Context, req *models.ListTransactionParams) (*ozon_cli.ListTransactionsResponse, error) {
	in := &ozon_cli.ListTransactionsParams{
		Filter: ozon_cli.ListTransactionsFilter{
			Date: ozon_cli.ListTransactionsFilterDate{
				From: req.Filter.Date.Value.From.Value,
				To:   req.Filter.Date.Value.To.Value,
			},
			OperationType:   req.Filter.OperationType,
			PostingNumber:   req.Filter.PostingNumber.Value,
			TransactionType: req.Filter.TransactionType.Value,
		},
		Page:     req.Page.Value,
		PageSize: req.PageSize.Value,
	}

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
				Name: models.NewOptString(item.Name),
				Sku:  models.NewOptInt64(item.SKU),
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
