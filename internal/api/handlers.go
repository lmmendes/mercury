package api

import (
	"mercury/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) createAccount(c echo.Context) error {
	ctx := c.Request().Context()

	var account models.Account
	if err := c.Bind(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := c.Validate(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.AccountService.Create(ctx, &account); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, account)
}

func (s *Server) getAccounts(c echo.Context) error {
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

	response, err := s.core.AccountService.List(ctx, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) getAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	account, err := s.core.AccountService.Get(c.Request().Context(), id)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if account == nil {
		return s.core.HandleError(nil, http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, account)
}

func (s *Server) updateAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var account models.Account
	if err := c.Bind(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}
	account.ID = id

	if err := c.Validate(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.AccountService.Update(c.Request().Context(), &account); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.core.AccountService.Delete(c.Request().Context(), id); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) createInbox(c echo.Context) error {
	accountID, _ := strconv.Atoi(c.Param("accountId"))
	var inbox models.Inbox
	if err := c.Bind(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}
	inbox.AccountID = accountID

	if err := c.Validate(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.InboxService.Create(c.Request().Context(), &inbox); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, inbox)
}

func (s *Server) getInboxes(c echo.Context) error {
	accountID, _ := strconv.Atoi(c.Param("accountId"))

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

	response, err := s.core.InboxService.ListByAccount(c.Request().Context(), accountID, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) getInbox(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	inbox, err := s.core.InboxService.Get(c.Request().Context(), inboxID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if inbox == nil {
		return s.core.HandleError(nil, http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, inbox)
}

func (s *Server) updateInbox(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	accountID, _ := strconv.Atoi(c.Param("accountId"))

	var inbox models.Inbox
	if err := c.Bind(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	// Set both ID and AccountID before validation
	inbox.ID = inboxID
	inbox.AccountID = accountID

	if err := c.Validate(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.InboxService.Update(c.Request().Context(), &inbox); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteInbox(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	if err := s.core.InboxService.Delete(c.Request().Context(), inboxID); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) createRule(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	var rule models.Rule
	if err := c.Bind(&rule); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}
	rule.InboxID = inboxID

	if err := c.Validate(&rule); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.RuleService.Create(c.Request().Context(), &rule); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, rule)
}

func (s *Server) getRules(c echo.Context) error {
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

	response, err := s.core.RuleService.ListByInbox(c.Request().Context(), inboxID, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) getRule(c echo.Context) error {
	ruleID, _ := strconv.Atoi(c.Param("ruleId"))
	rule, err := s.core.RuleService.Get(c.Request().Context(), ruleID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if rule == nil {
		return s.core.HandleError(nil, http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, rule)
}

func (s *Server) updateRule(c echo.Context) error {
	ruleID, _ := strconv.Atoi(c.Param("ruleId"))
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))

	var rule models.Rule
	if err := c.Bind(&rule); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	// Set both ID and InboxID before validation
	rule.ID = ruleID
	rule.InboxID = inboxID

	if err := c.Validate(&rule); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.RuleService.Update(c.Request().Context(), &rule); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteRule(c echo.Context) error {
	ruleID, _ := strconv.Atoi(c.Param("ruleId"))
	if err := s.core.RuleService.Delete(c.Request().Context(), ruleID); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

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
