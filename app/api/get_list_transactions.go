package api

import (
	models "SnipAndNeat/generated"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) ListTransactions(ctx *gin.Context) {
	req := &models.ListTransactionParams{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.log.Error("cannot parse request", zap.Error(err))
		HandleBadRequest(ctx, err)
		return
	}

	transactions, err := s.ozon.GetListTransaction(ctx, req)
	if err != nil {
		s.log.Error("cannot get transactions", zap.Error(err))
		HandleInternalError(ctx, err)
		return
	}

	HandleSuccess(ctx, transactions)
}
