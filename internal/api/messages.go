package api

import (
	"inbox451/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) getMessages(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))

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

	response, err := s.core.MessageService.ListByInbox(c.Request().Context(), inboxID, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) getMessage(c echo.Context) error {
	messageID, _ := strconv.Atoi(c.Param("messageId"))

	message, err := s.core.MessageService.Get(c.Request().Context(), messageID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if message == nil {
		return s.core.HandleError(nil, http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, message)
}
