package api

import (
	"SnipAndNeat/app/service/mailer"

	"github.com/diphantxm/ozon-api-client/ozon"
	"github.com/gin-gonic/gin"
)

type HealthInfo struct {
	// время работы
	Uptime string `json:"uptime"`
	// текущее время
	Time           string `json:"time"`
	Hostname       string `json:"hostname"`
	Application    string `json:"application"`
	Version        string `json:"version"`
	BuildTimestamp string `json:"build_timestamp"`
	BuildCommit    string `json:"build_commit"`
	BuildTag       string `json:"build_tag"`
}

// @Summary			Получение статуса
// @Description		Получение статуса
// @Tags			common
// @Produce			json
// @Param			X-Request-Id			header		string						false	"Случайная строка"
// @Success			200						{object}	HealthInfo					"Успешный запрос"
// @Failure			500						{object}	http.Error					"Ошибка"
// @Router			/health					[get]
func (s *Server) getHealth(ctx *gin.Context) {
	_, err := s.ozon.VientoProductsInfo(ctx, &ozon.GetListOfProductsParams{})
	if err != nil {
		HandleFailure(ctx, 500, err)
	}
	if err := s.mailer.SendText(&mailer.SmtpMessage{
		From:    s.cfg.EmailServer.Username,
		To:      []string{s.cfg.EmailServer.Recipient},
		Subject: "Test",
		Body:    "Test",
	}); err != nil {
		s.log.Sugar().Error(err)
	}

	s.log.Sugar().Info("get health info")
	HandleSuccess(ctx, nil)
}
