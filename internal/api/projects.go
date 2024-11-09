package api

import (
	"inbox451/internal/models"
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
