package api

import (
	"inbox451/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /users/:id/tokens/
func (s *Server) getTokens(c echo.Context) error {
	ctx := c.Request().Context()

	var query models.PaginationQuery
	if err := c.Bind(&query); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if err := c.Validate(&query); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	response, err := s.core.TokensService.List(ctx, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}
