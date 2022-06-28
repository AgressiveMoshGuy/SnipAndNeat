package api

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) VientoProducts(ctx *gin.Context) {
	products, err := s.ozon.ListProducts(ctx)
	if err != nil {
		HandleInternalError(ctx, err)
		return
	}

	HandleSuccess(ctx, products)
}
