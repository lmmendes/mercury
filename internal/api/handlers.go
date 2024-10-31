package api

import (
	"mercury/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) createAccount(c echo.Context) error {
	var account models.Account
	if err := c.Bind(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := c.Validate(&account); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.AccountService.Create(&account); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, account)
}

func (s *Server) getAccounts(c echo.Context) error {
	accounts, err := s.core.AccountService.List()
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, accounts)
}

func (s *Server) getAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	account, err := s.core.AccountService.Get(id)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if account == nil {
		return s.core.HandleError(err, http.StatusNotFound)
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

	if err := s.core.AccountService.Update(&account); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.core.AccountService.Delete(id); err != nil {
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

	if err := s.core.InboxService.Create(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, inbox)
}

func (s *Server) getInboxes(c echo.Context) error {
	accountID, _ := strconv.Atoi(c.Param("accountId"))
	inboxes, err := s.core.InboxService.GetByAccountID(accountID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, inboxes)
}

func (s *Server) getInbox(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	inbox, err := s.core.InboxService.Get(inboxID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if inbox == nil {
		return s.core.HandleError(err, http.StatusNotFound)
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

	if err := s.core.InboxService.Update(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteInbox(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	if err := s.core.InboxService.Delete(inboxID); err != nil {
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

	if err := s.core.RuleService.Create(&rule); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, rule)
}

func (s *Server) getRules(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	rules, err := s.core.RuleService.GetByInboxID(inboxID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, rules)
}

func (s *Server) getRule(c echo.Context) error {
	ruleID, _ := strconv.Atoi(c.Param("ruleId"))
	rule, err := s.core.RuleService.Get(ruleID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if rule == nil {
		return s.core.HandleError(err, http.StatusNotFound)
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

	if err := s.core.RuleService.Update(&rule); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteRule(c echo.Context) error {
	ruleID, _ := strconv.Atoi(c.Param("ruleId"))
	if err := s.core.RuleService.Delete(ruleID); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) getMessages(c echo.Context) error {
	inboxID, _ := strconv.Atoi(c.Param("inboxId"))
	messages, err := s.core.MessageService.GetByInboxID(inboxID)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, messages)
}
