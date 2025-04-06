package ozon

import (
	models "SnipAndNeat/generated"
	"context"
	"time"

	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

// GetSumTransactions возвращает агрегированную информацию о транзакциях за указанный период
func (z *OzonAPI) GetSumTransactions(ctx context.Context, in *models.GetSumProficiencyParamsDate) (*models.GetReportResponse, error) {
	if in == nil || in.From.Value.IsZero() || in.To.Value.IsZero() {
		return nil, ErrInvalidDateRange
	}

	// Проверяем, что дата начала не позже даты окончания
	if in.From.Value.After(in.To.Value) {
		return nil, ErrInvalidDateRange
	}

	// Получаем операции из БД
	operations, err := z.db.GetSumOperationsFromTime(ctx, in.From.Value, in.To.Value)
	if err != nil {
		z.log.Error().Err(err).
			Time("from", in.From.Value).
			Time("to", in.To.Value).
			Msg("failed to get operations from database")
		return nil, err
	}

	// Получаем транзакции из API Ozon
	transactions, err := z.client.Finance().GetTotalTransactionsSum(ctx, &ozon_cli.GetTotalTransactionsSumParams{
		TransactionType: "all",
		Date: ozon_cli.GetTotalTransactionsSumDate{
			From: in.From.Value,
			To:   in.To.Value,
		},
	})
	if err != nil {
		z.log.Error().Err(err).
			Time("from", in.From.Value).
			Time("to", in.To.Value).
			Msg("failed to get transactions from Ozon API")
		return nil, err
	}

	// Формируем отчет
	out := &models.GetReportResponse{
		Result: make([]models.ReportSumInfo, 0, len(operations)),
	}

	// Агрегируем данные по операциям
	for _, op := range operations {
		reportInfo := models.ReportSumInfo{
			Date:           op.Date,
			TransactionID:  op.TransactionID,
			OperationType:  string(op.OperationType),
			Amount:         op.Price,
			Commission:     op.Commission,
			ShippingCost:   op.ShippingCost,
			ReturnsCost:    op.ReturnsCost,
			MarketingCost:  op.MarketingCost,
			CategoryID:     op.CategoryID,
			CategoryName:   op.CategoryName,
			ItemsCount:     int64(op.ItemsCount),
			Profit:         op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost),
		}

		// Рассчитываем процент прибыли
		if reportInfo.Amount > 0 {
			reportInfo.ProfitPercentage = (reportInfo.Profit / reportInfo.Amount) * 100
		}

		out.Result = append(out.Result, reportInfo)
	}

	return out, nil
}

// GetDailyTransactions возвращает транзакции, сгруппированные по дням
func (z *OzonAPI) GetDailyTransactions(ctx context.Context, from, to time.Time) (map[time.Time]*models.ReportSumInfo, error) {
	operations, err := z.db.GetSumOperationsFromTime(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Группируем операции по дням
	dailyOps := make(map[time.Time]*models.ReportSumInfo)

	for _, op := range operations {
		// Округляем дату до начала дня
		day := time.Date(op.Date.Year(), op.Date.Month(), op.Date.Day(), 0, 0, 0, 0, op.Date.Location())
		
		if daily, exists := dailyOps[day]; exists {
			daily.Amount += op.Price
			daily.Commission += op.Commission
			daily.ShippingCost += op.ShippingCost
			daily.ReturnsCost += op.ReturnsCost
			daily.MarketingCost += op.MarketingCost
			daily.ItemsCount += int64(op.ItemsCount)
			daily.Profit = daily.Amount - (daily.Commission + daily.ShippingCost + daily.ReturnsCost + daily.MarketingCost)
			if daily.Amount > 0 {
				daily.ProfitPercentage = (daily.Profit / daily.Amount) * 100
			}
		} else {
			profit := op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost)
			dailyOps[day] = &models.ReportSumInfo{
				Date:           day,
				Amount:         op.Price,
				Commission:     op.Commission,
				ShippingCost:   op.ShippingCost,
				ReturnsCost:    op.ReturnsCost,
				MarketingCost:  op.MarketingCost,
				ItemsCount:     int64(op.ItemsCount),
				Profit:         profit,
				ProfitPercentage: 0,
			}
			if op.Price > 0 {
				dailyOps[day].ProfitPercentage = (profit / op.Price) * 100
			}
		}
	}

	return dailyOps, nil
}
