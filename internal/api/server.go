package api

import (
	"context"
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
	e.HideBanner = true
	s := &Server{
		core: core,
		echo: e,
	}

	// Add timeout middleware with a 30-second timeout
	e.Use(TimeoutMiddleware(30 * time.Second))

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Add middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	// Set custom error handler
	e.HTTPErrorHandler = s.errorHandler

	s.routes()
	return s
}

// Add the error handler method
func (s *Server) errorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
		if err := c.JSON(he.Code, he.Message); err != nil {
			s.core.Logger.Error("Failed to send error response: %v", err)
		}
		return
	}

	if err := s.core.HandleError(err, 500); err != nil {
		if err := c.JSON(500, err); err != nil {
			s.core.Logger.Error("Failed to send error response: %v", err)
		}
	}
}

func (s *Server) routes() {
	// Account routes
	s.echo.POST("/projects", s.createProject)
	s.echo.GET("/projects", s.getProjects)
	s.echo.GET("/projects/:id", s.getProject)
	s.echo.PUT("/projects/:id", s.updateProject)
	s.echo.DELETE("/projects/:id", s.deleteProject)

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
	return s.echo.Start(s.core.Config.Server.HTTP.Port)
}

// Add Shutdown method to Server struct
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
