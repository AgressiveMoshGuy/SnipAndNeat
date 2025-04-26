package ozon

import (
	"time"
)

// ProfitMetrics содержит расширенные метрики прибыльности
type ProfitMetrics struct {
	TotalRevenue     float64
	TotalCost        float64
	GrossProfit      float64
	NetProfit        float64
	ROI              float64
	MarginPercentage float64
	CommissionCost   float64
	ShippingCost     float64
	ReturnsCost      float64
	MarketingCost    float64
	// Новые метрики
	AverageProfitPerOrder float64
	ProfitByDay           map[time.Time]float64
	ProfitByCategory      map[string]float64
	ProfitTrend           []float64
	SeasonalFactors       map[int]float64 // Месяц -> коэффициент сезонности
}

// CategoryProfit содержит прибыльность по категориям
type CategoryProfit struct {
	CategoryID       string
	CategoryName     string
	Revenue          float64
	Cost             float64
	Profit           float64
	MarginPercentage float64
	ItemsCount       int
	// Новые метрики
	AverageProfitPerItem float64
	ReturnsRate          float64
	ProfitTrend          []float64
	SeasonalFactor       float64
}

// // GetDetailedProfitAnalysis возвращает детальный анализ прибыльности
// func (z *OzonAPI) GetDetailedProfitAnalysis(ctx context.Context, from, to time.Time) (*ProfitMetrics, error) {
// 	// Получаем все транзакции за период
// 	transactions, err := z.client.Finance().GetTotalTransactionsSum(ctx, &ozon_cli.GetTotalTransactionsSumParams{
// 		TransactionType: "all",
// 		Date: ozon_cli.GetTotalTransactionsSumDate{
// 			From: from,
// 			To:   to,
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Получаем операции из БД
// 	operations, err := z.db.GetSumOperationsFromTime(ctx, from, to)
// 	if err != nil {
// 		return nil, err
// 	}

// 	metrics := &ProfitMetrics{
// 		ProfitByDay:      make(map[time.Time]float64),
// 		ProfitByCategory: make(map[string]float64),
// 		SeasonalFactors:  make(map[int]float64),
// 	}

// 	// Расчет основных метрик
// 	orderCount := 0
// 	for _, op := range operations {
// 		// Округляем дату до начала дня
// 		day := time.Date(op.Date.Year(), op.Date.Month(), op.Date.Day(), 0, 0, 0, 0, op.Date.Location())

// 		// Рассчитываем прибыль для операции
// 		operationProfit := op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost)

// 		// Обновляем общие метрики
// 		metrics.TotalRevenue += op.Price
// 		metrics.CommissionCost += op.Commission
// 		metrics.ShippingCost += op.ShippingCost
// 		metrics.ReturnsCost += op.ReturnsCost
// 		metrics.MarketingCost += op.MarketingCost

// 		// Обновляем метрики по дням
// 		metrics.ProfitByDay[day] += operationProfit

// 		// Обновляем метрики по категориям
// 		metrics.ProfitByCategory[op.CategoryID] += operationProfit

// 		// Обновляем сезонные факторы
// 		month := int(op.Date.Month())
// 		metrics.SeasonalFactors[month] += operationProfit

// 		// Считаем количество заказов
// 		if op.OperationType != "return" {
// 			orderCount++
// 		}
// 	}

// 	// Расчет прибыли
// 	metrics.TotalCost = metrics.CommissionCost + metrics.ShippingCost + metrics.ReturnsCost + metrics.MarketingCost
// 	metrics.GrossProfit = metrics.TotalRevenue - metrics.TotalCost
// 	metrics.NetProfit = metrics.GrossProfit // Можно добавить дополнительные вычеты

// 	// Расчет ROI и маржи
// 	if metrics.TotalCost > 0 {
// 		metrics.ROI = (metrics.NetProfit / metrics.TotalCost) * 100
// 	}
// 	if metrics.TotalRevenue > 0 {
// 		metrics.MarginPercentage = (metrics.NetProfit / metrics.TotalRevenue) * 100
// 	}

// 	// Расчет средней прибыли на заказ
// 	if orderCount > 0 {
// 		metrics.AverageProfitPerOrder = metrics.NetProfit / float64(orderCount)
// 	}

// 	// Нормализация сезонных факторов
// 	totalProfit := 0.0
// 	for _, profit := range metrics.SeasonalFactors {
// 		totalProfit += profit
// 	}

// 	if totalProfit > 0 {
// 		for month, profit := range metrics.SeasonalFactors {
// 			metrics.SeasonalFactors[month] = profit / totalProfit * 12 // Нормализуем до годового распределения
// 		}
// 	}

// 	// Расчет тренда прибыли
// 	metrics.ProfitTrend = calculateProfitTrend(metrics.ProfitByDay)

// 	return metrics, nil
// }

// // GetCategoryProfitAnalysis возвращает анализ прибыльности по категориям
// func (z *OzonAPI) GetCategoryProfitAnalysis(ctx context.Context, from, to time.Time) ([]CategoryProfit, error) {
// 	operations, err := z.db.GetSumOperationsFromTime(ctx, from, to)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Группируем операции по категориям
// 	categoryMap := make(map[string]*CategoryProfit)

// 	for _, op := range operations {
// 		category, exists := categoryMap[op.CategoryID]
// 		if !exists {
// 			category = &CategoryProfit{
// 				CategoryID:   op.CategoryID,
// 				CategoryName: op.CategoryName,
// 			}
// 			categoryMap[op.CategoryID] = category
// 		}

// 		category.Revenue += op.Price
// 		category.Cost += op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost
// 		category.ItemsCount++

// 		// Подсчитываем возвраты
// 		if op.OperationType == "return" {
// 			category.ReturnsRate++
// 		}
// 	}

// 	// Преобразуем map в slice
// 	var result []CategoryProfit
// 	for _, category := range categoryMap {
// 		category.Profit = category.Revenue - category.Cost
// 		if category.Revenue > 0 {
// 			category.MarginPercentage = (category.Profit / category.Revenue) * 100
// 		}

// 		// Расчет средней прибыли на товар
// 		if category.ItemsCount > 0 {
// 			category.AverageProfitPerItem = category.Profit / float64(category.ItemsCount)
// 		}

// 		// Расчет процента возвратов
// 		if category.ItemsCount > 0 {
// 			category.ReturnsRate = category.ReturnsRate / float64(category.ItemsCount) * 100
// 		}

// 		// Расчет тренда прибыли по категории
// 		category.ProfitTrend = calculateCategoryProfitTrend(operations, category.CategoryID)

// 		// Расчет сезонного фактора
// 		category.SeasonalFactor = calculateSeasonalFactor(operations, category.CategoryID)

// 		result = append(result, *category)
// 	}

// 	return result, nil
// }

// // GetProfitForecast возвращает прогноз прибыли на основе исторических данных
// func (z *OzonAPI) GetProfitForecast(ctx context.Context, days int) (*ProfitMetrics, error) {
// 	// Получаем исторические данные за последние 90 дней для лучшего прогноза
// 	from := time.Now().AddDate(0, 0, -90)
// 	to := time.Now()

// 	historicalMetrics, err := z.GetDetailedProfitAnalysis(ctx, from, to)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Рассчитываем среднюю дневную прибыль с учетом сезонности
// 	currentMonth := int(time.Now().Month())
// 	seasonalFactor := historicalMetrics.SeasonalFactors[currentMonth]

// 	// Если сезонный фактор не рассчитан, используем среднее значение
// 	if seasonalFactor == 0 {
// 		seasonalFactor = 1.0
// 	}

// 	// Простой прогноз на основе среднего значения и сезонности
// 	avgDailyProfit := historicalMetrics.NetProfit / 90.0
// 	forecast := &ProfitMetrics{
// 		TotalRevenue:          historicalMetrics.TotalRevenue * float64(days) / 90 * seasonalFactor,
// 		TotalCost:             historicalMetrics.TotalCost * float64(days) / 90 * seasonalFactor,
// 		GrossProfit:           historicalMetrics.GrossProfit * float64(days) / 90 * seasonalFactor,
// 		NetProfit:             historicalMetrics.NetProfit * float64(days) / 90 * seasonalFactor,
// 		ROI:                   historicalMetrics.ROI,
// 		MarginPercentage:      historicalMetrics.MarginPercentage,
// 		CommissionCost:        historicalMetrics.CommissionCost * float64(days) / 90 * seasonalFactor,
// 		ShippingCost:          historicalMetrics.ShippingCost * float64(days) / 90 * seasonalFactor,
// 		ReturnsCost:           historicalMetrics.ReturnsCost * float64(days) / 90 * seasonalFactor,
// 		MarketingCost:         historicalMetrics.MarketingCost * float64(days) / 90 * seasonalFactor,
// 		AverageProfitPerOrder: historicalMetrics.AverageProfitPerOrder,
// 		ProfitByDay:           make(map[time.Time]float64),
// 		ProfitByCategory:      historicalMetrics.ProfitByCategory,
// 		ProfitTrend:           historicalMetrics.ProfitTrend,
// 		SeasonalFactors:       historicalMetrics.SeasonalFactors,
// 	}

// 	// Заполняем прогноз по дням
// 	startDate := time.Now()
// 	for i := 0; i < days; i++ {
// 		date := startDate.AddDate(0, 0, i)
// 		month := int(date.Month())
// 		monthFactor := historicalMetrics.SeasonalFactors[month]
// 		if monthFactor == 0 {
// 			monthFactor = 1.0
// 		}
// 		forecast.ProfitByDay[date] = avgDailyProfit * monthFactor
// 	}

// 	return forecast, nil
// }

// // Вспомогательные функции

// // calculateProfitTrend рассчитывает тренд прибыли
// func calculateProfitTrend(profitByDay map[time.Time]float64) []float64 {
// 	// Сортируем даты
// 	var dates []time.Time
// 	for date := range profitByDay {
// 		dates = append(dates, date)
// 	}

// 	// Простой расчет тренда (можно заменить на более сложный алгоритм)
// 	var trend []float64
// 	if len(dates) > 0 {
// 		// Используем скользящее среднее для сглаживания
// 		windowSize := 7 // Недельное окно
// 		for i := 0; i < len(dates); i++ {
// 			sum := 0.0
// 			count := 0

// 			// Считаем среднее за окно
// 			for j := max(0, i-windowSize); j < min(len(dates), i+windowSize+1); j++ {
// 				sum += profitByDay[dates[j]]
// 				count++
// 			}

// 			if count > 0 {
// 				trend = append(trend, sum/float64(count))
// 			}
// 		}
// 	}

// 	return trend
// }

// // calculateCategoryProfitTrend рассчитывает тренд прибыли по категории
// func calculateCategoryProfitTrend(operations []ozon_cli.ListTransactionsResultOperation, categoryID string) []float64 {
// 	// Группируем операции по дням
// 	profitByDay := make(map[time.Time]float64)

// 	for _, op := range operations {
// 		if op.CategoryID == categoryID {
// 			day := time.Date(op.Date.Year(), op.Date.Month(), op.Date.Day(), 0, 0, 0, 0, op.Date.Location())
// 			profit := op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost)
// 			profitByDay[day] += profit
// 		}
// 	}

// 	return calculateProfitTrend(profitByDay)
// }

// // calculateSeasonalFactor рассчитывает сезонный фактор для категории
// func calculateSeasonalFactor(operations []ozon_cli.ListTransactionsResultOperation, categoryID string) float64 {
// 	// Группируем операции по месяцам
// 	profitByMonth := make(map[int]float64)

// 	for _, op := range operations {
// 		if op.CategoryID == categoryID {
// 			month := int(op.Date.Month())
// 			profit := op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost)
// 			profitByMonth[month] += profit
// 		}
// 	}

// 	// Рассчитываем среднюю прибыль по месяцам
// 	totalProfit := 0.0
// 	monthCount := 0

// 	for _, profit := range profitByMonth {
// 		totalProfit += profit
// 		monthCount++
// 	}

// 	if monthCount == 0 {
// 		return 1.0
// 	}

// 	avgProfit := totalProfit / float64(monthCount)

// 	// Находим месяц с максимальной прибылью
// 	maxProfit := 0.0
// 	for _, profit := range profitByMonth {
// 		if profit > maxProfit {
// 			maxProfit = profit
// 		}
// 	}

// 	if avgProfit > 0 {
// 		return maxProfit / avgProfit
// 	}

// 	return 1.0
// }

// // Вспомогательные функции для min/max
// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }
