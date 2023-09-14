package usecase

import "github.com/anti-duhring/easy-shipping-microsservice/internal/entity"

type ChangeRouteStatusInput struct {
	ID         string
	StartedAt  entity.CustomTime
	FinishedAt entity.CustomTime
	Event      string
}
