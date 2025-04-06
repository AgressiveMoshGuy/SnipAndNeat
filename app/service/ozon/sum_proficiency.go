package ozon

import (
	models "SnipAndNeat/generated"
	"context"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

func (s *OzonAPI) GetSumProficiency(ctx context.Context, req *models.GetSumProficiencyParams) error {
	_, err := s.client.Finance().GetTotalTransactionsSum(ctx, &ozon_cli.GetTotalTransactionsSumParams{
		TransactionType: "all",
		Date: ozon_cli.GetTotalTransactionsSumDate{
			From: req.Date.Value.From.Value,
			To:   req.Date.Value.To.Value,
		},
	})

	return err
}
