package ozon

import (
	"context"
	"time"

	models "SnipAndNeat/generated"
	ozon_cli "github.com/diphantxm/ozon-api-client/ozon"
)

// SalesMetrics содержит метрики продаж
type SalesMetrics struct {
	TotalSales       float64
	TotalOrders      int
	AverageOrderValue float64
	ItemsSold        int
	ReturnsCount     int
	ReturnsRate      float64
	SalesByDay       map[time.Time]float64
	OrdersByDay      map[time.Time]int
}

// ProductPerformance содержит метрики производительности товара
type ProductPerformance struct {
	ProductID        string
	ProductName      string
	SKU              string
	SalesCount       int
	Revenue          float64
	Profit           float64
	MarginPercentage float64
	ReturnsCount     int
	ReturnsRate      float64
	Views            int
	ConversionRate   float64
}

// GetSalesAnalytics возвращает детальный анализ продаж
func (z *OzonAPI) GetSalesAnalytics(ctx context.Context, from, to time.Time) (*SalesMetrics, error) {
	// Получаем транзакции из API Ozon
	transactions, err := z.client.Finance().GetTotalTransactionsSum(ctx, &ozon_cli.GetTotalTransactionsSumParams{
		TransactionType: "all",
		Date: ozon_cli.GetTotalTransactionsSumDate{
			From: from,
			To:   to,
		},
	})
	if err != nil {
		return nil, err
	}

	// Получаем операции из БД
	operations, err := z.db.GetSumOperationsFromTime(ctx, from, to)
	if err != nil {
		return nil, err
	}

	metrics := &SalesMetrics{
		SalesByDay:  make(map[time.Time]float64),
		OrdersByDay: make(map[time.Time]int),
	}

	// Анализируем операции
	for _, op := range operations {
		// Округляем дату до начала дня
		day := time.Date(op.Date.Year(), op.Date.Month(), op.Date.Day(), 0, 0, 0, 0, op.Date.Location())
		
		// Обновляем метрики
		metrics.TotalSales += op.Price
		metrics.TotalOrders++
		metrics.ItemsSold += op.ItemsCount
		
		// Обновляем метрики по дням
		metrics.SalesByDay[day] += op.Price
		metrics.OrdersByDay[day]++
		
		// Подсчитываем возвраты
		if op.OperationType == "return" {
			metrics.ReturnsCount++
		}
	}

	// Рассчитываем средний чек
	if metrics.TotalOrders > 0 {
		metrics.AverageOrderValue = metrics.TotalSales / float64(metrics.TotalOrders)
	}

	// Рассчитываем процент возвратов
	if metrics.TotalOrders > 0 {
		metrics.ReturnsRate = float64(metrics.ReturnsCount) / float64(metrics.TotalOrders) * 100
	}

	return metrics, nil
}

// GetProductPerformance возвращает анализ производительности товаров
func (z *OzonAPI) GetProductPerformance(ctx context.Context, from, to time.Time) ([]ProductPerformance, error) {
	// Получаем список товаров
	products, err := z.client.Products().GetListOfProducts(ctx, &ozon_cli.GetListOfProductsParams{})
	if err != nil {
		return nil, err
	}

	// Получаем операции из БД
	operations, err := z.db.GetSumOperationsFromTime(ctx, from, to)
	if err != nil {
		return nil, err
	}

	// Группируем операции по товарам
	productMap := make(map[string]*ProductPerformance)

	// Инициализируем карту товаров
	for _, product := range products.Result.Items {
		productMap[product.ProductID] = &ProductPerformance{
			ProductID:   product.ProductID,
			ProductName: product.Name,
			SKU:         product.SKU,
		}
	}

	// Анализируем операции
	for _, op := range operations {
		product, exists := productMap[op.ProductID]
		if !exists {
			continue
		}

		product.SalesCount++
		product.Revenue += op.Price
		product.Profit += op.Price - (op.Commission + op.ShippingCost + op.ReturnsCost + op.MarketingCost)

		if op.OperationType == "return" {
			product.ReturnsCount++
		}
	}

	// Рассчитываем дополнительные метрики
	var result []ProductPerformance
	for _, product := range productMap {
		if product.SalesCount > 0 {
			product.MarginPercentage = (product.Profit / product.Revenue) * 100
			product.ReturnsRate = float64(product.ReturnsCount) / float64(product.SalesCount) * 100
		}
		result = append(result, *product)
	}

	return result, nil
}

// GetSalesComparison сравнивает продажи с предыдущим периодом
func (z *OzonAPI) GetSalesComparison(ctx context.Context, from, to time.Time) (map[string]float64, error) {
	// Рассчитываем длительность периода
	duration := to.Sub(from)
	
	// Получаем данные за текущий период
	currentMetrics, err := z.GetSalesAnalytics(ctx, from, to)
	if err != nil {
		return nil, err
	}
	
	// Получаем данные за предыдущий период
	prevFrom := from.Add(-duration)
	prevTo := from
	prevMetrics, err := z.GetSalesAnalytics(ctx, prevFrom, prevTo)
	if err != nil {
		return nil, err
	}
	
	// Рассчитываем изменения
	comparison := make(map[string]float64)
	
	// Изменение общей выручки
	if prevMetrics.TotalSales > 0 {
		comparison["revenue_change"] = (currentMetrics.TotalSales - prevMetrics.TotalSales) / prevMetrics.TotalSales * 100
	}
	
	// Изменение количества заказов
	if prevMetrics.TotalOrders > 0 {
		comparison["orders_change"] = float64(currentMetrics.TotalOrders-prevMetrics.TotalOrders) / float64(prevMetrics.TotalOrders) * 100
	}
	
	// Изменение среднего чека
	if prevMetrics.AverageOrderValue > 0 {
		comparison["avg_order_change"] = (currentMetrics.AverageOrderValue - prevMetrics.AverageOrderValue) / prevMetrics.AverageOrderValue * 100
	}
	
	// Изменение количества проданных товаров
	if prevMetrics.ItemsSold > 0 {
		comparison["items_sold_change"] = float64(currentMetrics.ItemsSold-prevMetrics.ItemsSold) / float64(prevMetrics.ItemsSold) * 100
	}
	
	// Изменение процента возвратов
	comparison["returns_rate_change"] = currentMetrics.ReturnsRate - prevMetrics.ReturnsRate
	
	return comparison, nil
} 