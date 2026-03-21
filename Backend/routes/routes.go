package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vihantandon/Coders_Hub/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

}
