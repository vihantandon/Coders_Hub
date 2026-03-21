package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
)

func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := boot.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := boot.DB.Where("email = ?", input.Email).First(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Email"})
		return
	}

	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
