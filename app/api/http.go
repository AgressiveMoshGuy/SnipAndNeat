package api

import (
	"SnipAndNeat/app/config"
	"SnipAndNeat/app/service/mailer"
	"SnipAndNeat/app/service/ozon"
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func New(cfg *config.Config, mailer mailer.EmailSender, ozon ozon.Ozon) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	logger, _ := zap.NewDevelopment()

	s := Server{
		log:    logger.Named("http"),
		cfg:    cfg,
		rtr:    router,
		mailer: mailer,
		ozon:   ozon,
		srv: &http.Server{
			Addr:           cfg.Server.Address,
			Handler:        router,
			ReadTimeout:    cfg.Server.ReadTimeout,
			WriteTimeout:   cfg.Server.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		// validator: validator.New(),
	}

	//for name, validator := range map[string]validator.Func{
	//	"phone":           validatePhone,
	//	"name":            validateName,
	//	"birth_place":     validateBirthPlace,
	//	"date":            validateDate,
	//	"issured_by_code": validateIssuredByCode,
	//	"issured_by":      validateIssuedBy,
	//} {
	//	err := s.validator.RegisterValidation(name, validator)
	//	if err != nil {
	//		s.log.Error().Err(err).Msg("cannot register validation")
	//		return nil, err
	//	}
	//}

	return &s, nil
}

type Server struct {
	log    *zap.Logger
	cfg    *config.Config
	rtr    *gin.Engine
	srv    *http.Server
	mailer mailer.EmailSender
	ozon   ozon.Ozon
	// validator *validator.Validate
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) Start(context.Context) error {
	// создаем метрику, которая будет отображать сумму продаж
	// для этого нужно создать метрику типа GaugeVec
	// GaugeVec - это тип метрики, который может иметь несколько значений
	// в этом случае мы будем иметь метрику, которая будет отображать сумму
	// продаж по каждому типу товара
	// для этого нужно создать метрику с именем, например "sells_sum"
	// и добавить лэйбл "type", который будет отображать тип товара
	// после этого мы можем использовать методы Inc и Dec для увеличения
	// или уменьшения значения метрики
	// например, если у нас есть заказ на товар "apple", то мы можем
	// использовать следующий код:
	// metric.WithLabelValues("apple").Inc()
	// это увеличит значение метрики "sells_sum" для типа "apple" на 1
	// мы можем использовать этот код в любой момент, когда мы хотим
	// отобразить сумму продаж по типу товара
	/*************  ✨ Codeium Command 🌟  *************/
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sells_sum",
			Help: "Sum of sells by type",
		},
		[]string{"type"},
	)
	prometheus.MustRegister(metric)

	metric.WithLabelValues("apple").Inc()
	s.rtr.Handle("GET", "/metrics", prometheusHandler())

	s.rtr.Group("/").
		GET("/", s.getHealth).
		GET("/health", s.getHealth).
		GET("/healthz", s.getHealth).
		GET("/health-check", s.getHealth)

	s.rtr.Group("/ozon").
		POST("/list_transactions", s.ListTransactions).
		POST("/sum/services", s.GetSumServices)
	s.rtr.Group("/viento").
		POST("/list_products", s.VientoProducts).
		POST("/update_barcodes", s.UpdateItemsBarcode)
	// GET("/dictionary", s.getDictionary).
	// GET("/credit_products", s.getCreditProducts).
	// GET("/draft", s.getDraft).
	// PUT("/draft", s.putDraft).
	// DELETE("/draft", s.deleteDraft)

	//s.rtr.Group("/code").Use(s.wrap).Use(s.auth).
	//	POST("/send", s.postCodeSend).
	//	POST("/verify", s.postCodeVerify).
	//	POST("/validate", s.postCodeValidate)

	s.log.Sugar().Debugf("start listening %q", s.cfg.Server.Address)
	errCh := make(chan error)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			errCh <- errors.Wrap(err, "cannot listen and serve")
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(s.cfg.Server.StartTimeout):
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error { return s.srv.Shutdown(ctx) }

func HandleSuccess(ctx *gin.Context, res interface{}) {
	ctx.JSON(http.StatusOK, res)
}

func HandleBadRequest(ctx *gin.Context, err error) {
	HandleFailure(ctx, http.StatusBadRequest, err)
}

func HandleForbidden(ctx *gin.Context, err error) {
	HandleFailure(ctx, http.StatusForbidden, err)
}

func HandleNotFound(ctx *gin.Context, err error) {
	HandleFailure(ctx, http.StatusNotFound, err)
}

func HandleInternalError(ctx *gin.Context, err error) {
	HandleFailure(ctx, http.StatusInternalServerError, err)
}

type Error struct {
	Error ErrorInfo `json:"error"`
}
type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleFailure(ctx *gin.Context, status int, err error) {
	if err != nil {
		ctx.JSON(status, Error{Error: ErrorInfo{Code: status, Message: err.Error()}})
	}
	ctx.Status(status)
}

type ResponseWriter struct {
	body *bytes.Buffer
	gin.ResponseWriter
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
