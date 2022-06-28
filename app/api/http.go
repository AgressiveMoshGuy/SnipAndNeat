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

func (s *Server) Start(context.Context) error {
	s.rtr.Group("/").
		GET("/", s.getHealth).
		GET("/health", s.getHealth).
		GET("/healthz", s.getHealth).
		GET("/health-check", s.getHealth)

	s.rtr.Group("/viento").
		POST("/list_products", s.VientoProducts)
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
