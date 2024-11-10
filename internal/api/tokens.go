package api

import (
	"inbox451/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GET /users/:userId/tokens/
func (s *Server) GetTokensByUser(c echo.Context) error {
	ctx := c.Request().Context()
	userId, _ := strconv.Atoi(c.Param("userId"))

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

	response, err := s.core.TokensService.ListByUser(ctx, userId, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

// GET /users/:userId/tokens/:tokenId
func (s *Server) GetTokenByUser(c echo.Context) error {
	userID, _ := strconv.Atoi(c.Param("userId"))
	tokenID, _ := strconv.Atoi(c.Param("tokenId"))

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

	response, err := s.core.TokensService.GetByUser(c.Request().Context(), tokenID, userID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

// POST /users/:userId/tokens/
func (s *Server) CreateTokenForUser(c echo.Context) error {
	ctx := c.Request().Context()

	input := new(models.Token)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&input); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	userId, _ := strconv.Atoi(c.Param("userId"))
	newToken, err := s.core.TokensService.CreateTokenForUser(ctx, userId, input)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, newToken)
}
