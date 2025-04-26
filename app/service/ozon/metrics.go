package ozon

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Метрики транзакций
	transactionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ozon_transaction_total",
			Help: "Общее количество транзакций по типам",
		},
		[]string{"type", "category"},
	)

	transactionAmount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ozon_transaction_amount",
			Help: "Сумма транзакций по типам",
		},
		[]string{"type", "category"},
	)

	// Метрики прибыли
	profitTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ozon_profit_total",
			Help: "Общая прибыль",
		},
	)

	profitByDay = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ozon_profit_by_day",
			Help: "Прибыль по дням",
		},
		[]string{"date"},
	)

	profitByCategory = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ozon_profit_by_category",
			Help: "Прибыль по категориям",
		},
		[]string{"category_id", "category_name"},
	)

	// Метрики возвратов
	returnsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ozon_returns_total",
			Help: "Общее количество возвратов",
		},
	)

	returnsAmount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ozon_returns_amount",
			Help: "Сумма возвратов",
		},
	)

	// Метрики комиссий
	commissionTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ozon_commission_total",
			Help: "Общая сумма комиссий",
		},
	)

	commissionByType = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ozon_commission_by_type",
			Help: "Комиссии по типам",
		},
		[]string{"type"},
	)
)

// UpdateTransactionMetrics обновляет метрики транзакций
func UpdateTransactionMetrics(operationType, categoryID, categoryName string, amount float64) {
	transactionTotal.WithLabelValues(operationType, categoryID).Inc()
	transactionAmount.WithLabelValues(operationType, categoryID).Add(amount)
}

// UpdateProfitMetrics обновляет метрики прибыли
func UpdateProfitMetrics(profit float64, date time.Time, categoryID, categoryName string) {
	profitTotal.Set(profit)
	profitByDay.WithLabelValues(date.Format("2006-01-02")).Set(profit)
	profitByCategory.WithLabelValues(categoryID, categoryName).Set(profit)
}

// UpdateReturnsMetrics обновляет метрики возвратов
func UpdateReturnsMetrics(amount float64) {
	returnsTotal.Inc()
	returnsAmount.Set(amount)
}

// UpdateCommissionMetrics обновляет метрики комиссий
func UpdateCommissionMetrics(amount float64, commissionType string) {
	commissionTotal.Set(amount)
	commissionByType.WithLabelValues(commissionType).Set(amount)
}

// ResetMetrics сбрасывает все метрики
func ResetMetrics() {
	transactionTotal.Reset()
	transactionAmount.Reset()
	profitTotal.Set(0)
	profitByDay.Reset()
	profitByCategory.Reset()
	returnsAmount.Set(0)
	commissionTotal.Set(0)
	commissionByType.Reset()
}
