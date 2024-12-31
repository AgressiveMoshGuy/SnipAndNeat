package api

import (
	models "SnipAndNeat/generated"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) GetSumServices(ctx *gin.Context) {
	req := &models.SumServicesByDayParams{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.log.Error("cannot parse request", zap.Error(err))
		HandleBadRequest(ctx, err)
		return
	}

	sum, err := s.ozon.GetSumServices(ctx, req)
	if err != nil {
		HandleInternalError(ctx, err)
		return
	}

	HandleSuccess(ctx, sum)
}
