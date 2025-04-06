package ozon

import "errors"

var (
	// ErrInvalidDateRange возвращается при некорректном диапазоне дат
	ErrInvalidDateRange = errors.New("invalid date range")

	// ErrNoData возвращается при отсутствии данных
	ErrNoData = errors.New("no data available")

	// ErrInvalidOperation возвращается при некорректной операции
	ErrInvalidOperation = errors.New("invalid operation")

	// ErrAPIError возвращается при ошибке API Ozon
	ErrAPIError = errors.New("ozon API error")

	// ErrDatabaseError возвращается при ошибке базы данных
	ErrDatabaseError = errors.New("database error")
) 