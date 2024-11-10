package api

import "github.com/labstack/echo/v4"

func (s *Server) routes(api *echo.Group) {

	// User routes
	api.POST("/users", s.createUser)
	api.GET("/users", s.getUsers)
	api.GET("/users/:userId", s.getUser)
	api.PUT("/users/:userId", s.updateUser)
	api.DELETE("/users/:userId", s.deleteUser)
	api.GET("/users/:userId/projects", s.getProjectsByUser)

	// ProjectUser routes
	api.POST("/projects/:projectId/users", s.projectAddUser)
	api.DELETE("/projects/:projectId/users/:userId", s.projectRemoveUser)

	// Project routes
	api.POST("/projects", s.createProject)
	api.GET("/projects", s.getProjects)
	api.GET("/projects/:projectId", s.getProject)
	api.PUT("/projects/:projectId", s.updateProject)
	api.DELETE("/projects/:projectId", s.deleteProject)

	// Token routes
	api.GET("/users/:userId/tokens", s.ListTokensByUser)
	api.GET("/users/:userId/tokens/:tokenId", s.GetTokenByUser)
	api.POST("/users/:userId/tokens", s.CreateTokenForUser)
	api.DELETE("/users/:userId/tokens/:tokenId", s.DeleteTokenByUser)

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
