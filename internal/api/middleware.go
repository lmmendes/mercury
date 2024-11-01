package api

import (
	"mercury/internal/core"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) errorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = core.APIError{
			Code:    code,
			Message: err.Error(),
		}
	}

	// Send error response
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			s.core.Logger.Error("Failed to send error response: %v", err)
		}
	}
}
