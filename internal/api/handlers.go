package api

import (
	"mercury/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) createProject(c echo.Context) error {
	ctx := c.Request().Context()

	var project models.Project
	if err := c.Bind(&project); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := c.Validate(&project); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.ProjectService.Create(ctx, &project); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, project)
}

func (s *Server) getProjects(c echo.Context) error {
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

	response, err := s.core.ProjectService.List(ctx, query.Limit, query.Offset)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) getProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	project, err := s.core.ProjectService.Get(c.Request().Context(), id)
	if err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	if project == nil {
		return s.core.HandleError(nil, http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, project)
}

func (s *Server) updateProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var project models.Project
	if err := c.Bind(&project); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}
	project.ID = id

	if err := c.Validate(&project); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.ProjectService.Update(c.Request().Context(), &project); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.core.ProjectService.Delete(c.Request().Context(), id); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) createInbox(c echo.Context) error {
	projectID, _ := strconv.Atoi(c.Param("projectId"))
	var inbox models.Inbox
	if err := c.Bind(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}
	inbox.ProjectID = projectID

	if err := c.Validate(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	if err := s.core.InboxService.Create(c.Request().Context(), &inbox); err != nil {
		return s.core.HandleError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, inbox)
}

func (s *Server) getInboxes(c echo.Context) error {
	projectID, _ := strconv.Atoi(c.Param("projectId"))

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

	response, err := s.core.InboxService.ListByProject(c.Request().Context(), projectID, query.Limit, query.Offset)
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
	projectID, _ := strconv.Atoi(c.Param("projectId"))

	var inbox models.Inbox
	if err := c.Bind(&inbox); err != nil {
		return s.core.HandleError(err, http.StatusBadRequest)
	}

	// Set both ID and ProjectID before validation
	inbox.ID = inboxID
	inbox.ProjectID = projectID

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
	var rule models.ForwardRule
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

	var rule models.ForwardRule
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
