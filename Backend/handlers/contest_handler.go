package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
)

func GetContests(c *gin.Context) {
	platform := c.Query("platform")

	var contests []models.Contest
	query := boot.DB.Order("start ASC")

	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	if err := query.Find(&contests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":    len(contests),
		"contests": contests,
	})
}
