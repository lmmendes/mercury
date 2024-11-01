package api

import (
	"mercury/internal/core"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type Server struct {
	core *core.Core
	echo *echo.Echo
}

func NewServer(core *core.Core) *Server {
	e := echo.New()
	s := &Server{
		core: core,
		echo: e,
	}

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	e.HTTPErrorHandler = s.errorHandler

	s.routes()
	return s
}

func (s *Server) routes() {
	// Account routes
	s.echo.POST("/accounts", s.createAccount)
	s.echo.GET("/accounts", s.getAccounts)
	s.echo.GET("/accounts/:id", s.getAccount)
	s.echo.PUT("/accounts/:id", s.updateAccount)
	s.echo.DELETE("/accounts/:id", s.deleteAccount)

	// Inbox routes
	s.echo.POST("/accounts/:accountId/inboxes", s.createInbox)
	s.echo.GET("/accounts/:accountId/inboxes", s.getInboxes)
	s.echo.GET("/accounts/:accountId/inboxes/:inboxId", s.getInbox)
	s.echo.PUT("/accounts/:accountId/inboxes/:inboxId", s.updateInbox)
	s.echo.DELETE("/accounts/:accountId/inboxes/:inboxId", s.deleteInbox)

	// Rule routes
	s.echo.POST("/accounts/:accountId/inboxes/:inboxId/rules", s.createRule)
	s.echo.GET("/accounts/:accountId/inboxes/:inboxId/rules", s.getRules)
	s.echo.GET("/accounts/:accountId/inboxes/:inboxId/rules/:ruleId", s.getRule)
	s.echo.PUT("/accounts/:accountId/inboxes/:inboxId/rules/:ruleId", s.updateRule)
	s.echo.DELETE("/accounts/:accountId/inboxes/:inboxId/rules/:ruleId", s.deleteRule)

	// Message routes
	s.echo.GET("/accounts/:accountId/inboxes/:inboxId/messages", s.getMessages)
}

func (s *Server) ListenAndServe() error {
	return s.echo.Start(s.core.Config.HTTPPort)
}
