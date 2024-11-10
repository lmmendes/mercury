package api

import "github.com/labstack/echo/v4"

func (s *Server) routes(api *echo.Group) {

	// Token routes
	api.GET("/users/:userId/tokens", s.ListTokensByUser)
	api.GET("/users/:userId/tokens/:tokenId", s.GetTokenByUser)
	api.POST("/users/:userId/tokens", s.CreateTokenForUser)
	api.DELETE("/users/:userId/tokens/:tokenId", s.DeleteTokenByUser)

	// Project routes
	api.POST("/projects", s.createProject)
	api.GET("/projects", s.getProjects)
	api.GET("/projects/:id", s.getProject)
	api.PUT("/projects/:id", s.updateProject)
	api.DELETE("/projects/:id", s.deleteProject)

	// Inbox routes
	api.POST("/projects/:projectId/inboxes", s.createInbox)
	api.GET("/projects/:projectId/inboxes", s.getInboxes)
	api.GET("/projects/:projectId/inboxes/:inboxId", s.getInbox)
	api.PUT("/projects/:projectId/inboxes/:inboxId", s.updateInbox)
	api.DELETE("/projects/:projectId/inboxes/:inboxId", s.deleteInbox)

	// Rule routes
	api.POST("/projects/:projectId/inboxes/:inboxId/rules", s.createRule)
	api.GET("/projects/:projectId/inboxes/:inboxId/rules", s.getRules)
	api.GET("/projects/:projectId/inboxes/:inboxId/rules/:ruleId", s.getRule)
	api.PUT("/projects/:projectId/inboxes/:inboxId/rules/:ruleId", s.updateRule)
	api.DELETE("/projects/:projectId/inboxes/:inboxId/rules/:ruleId", s.deleteRule)

	// Message routes
	api.GET("/projects/:projectId/inboxes/:inboxId/messages", s.getMessages)
}
