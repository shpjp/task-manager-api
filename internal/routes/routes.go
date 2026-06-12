package routes

import (
	"net/http"

	"task-manager-api/internal/auth"
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	Tasks       *handlers.TaskHandler
	Attachments *handlers.AttachmentHandler
	Events      *handlers.EventsHandler
}

func RegisterRoutes(r *gin.Engine, h Handlers, jwt *auth.JWTManager) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", h.Auth.Signup)
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/logout", h.Auth.Logout)
		authGroup.GET("/me", middleware.RequireAuth(jwt), h.Auth.Me)
	}

	tasks := r.Group("/tasks", middleware.RequireAuth(jwt))
	{
		tasks.POST("", h.Tasks.Create)
		tasks.GET("", h.Tasks.List)
		tasks.GET("/:id", h.Tasks.Get)
		tasks.PATCH("/:id", h.Tasks.Update)
		tasks.DELETE("/:id", h.Tasks.Delete)

		tasks.GET("/:id/activity", h.Tasks.Activity)

		tasks.POST("/:id/attachments", h.Attachments.Upload)
		tasks.GET("/:id/attachments", h.Attachments.List)
		tasks.DELETE("/:id/attachments/:attachmentID", h.Attachments.Delete)
	}

	r.GET("/events", middleware.RequireAuth(jwt), h.Events.Stream)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{"code": "NOT_FOUND", "message": "Route not found"},
		})
	})
}
