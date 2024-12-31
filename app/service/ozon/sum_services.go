package ozon

import (
	models "SnipAndNeat/generated"
	"context"
	"time"
)

func (s *OzonAPI) GetSumServices(ctx context.Context, in *models.SumServicesByDayParams) (*models.GetSumServicesByDayOKApplicationJSON, error) {
	sumServices, err := s.db.GetSumServices(ctx, time.Time(in.Date.Value))
	if err != nil {
		return nil, err
	}

	result := models.GetSumServicesByDayOKApplicationJSON{}

	for operation_name, amount := range sumServices {
		result = append(result, models.GetSumServicesByDayOKItem{
			Name:   operation_name,
			Amount: float32(amount),
		})
	}

	return &result, nil
}
