package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	CommitSHA string `json:"commitSha"`
}

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Version:   s.core.Version,
		CommitSHA: s.core.Commit,
	})
}
