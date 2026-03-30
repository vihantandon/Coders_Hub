package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vihantandon/Coders_Hub/handlers"
	"github.com/vihantandon/Coders_Hub/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/contests", handlers.GetContests)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/reminders", handlers.CreateReminder)
		protected.GET("/reminders", handlers.GetReminders)
		protected.DELETE("/reminders/:id", handlers.DeleteReminder)
	}
}
