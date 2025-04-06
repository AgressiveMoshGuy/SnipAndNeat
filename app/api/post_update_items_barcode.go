package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) UpdateItemsBarcode(ctx *gin.Context) {
	res, err := s.ozon.UpdateEANCodesWithItems(ctx)
	if err != nil {
		HandleInternalError(ctx, err)
		return
	}

	HandleSuccess(ctx, res)
}
