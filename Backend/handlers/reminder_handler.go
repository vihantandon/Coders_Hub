package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
)

type CreateReminderInput struct {
	ContestID   uint `json:"contest_id" binding:"required"`
	HoursBefore uint `json:"hours_before" binding:"required,min=1,max=168"`
}

// POST /reminders
func CreateReminder(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	var input CreateReminderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var contest models.Contest
	if err := boot.DB.First(&contest, input.ContestID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05", contest.Start)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid contest start time format"})
		return
	}

	if startTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contest has already started"})
		return
	}

	sendAt := startTime.Add(-time.Duration(input.HoursBefore) * time.Hour)

	if sendAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reminder time is already in the past. Choose fewer hours before the contest.",
		})
		return
	}

	var existing models.Reminder
	if err := boot.DB.Where("user_id = ? AND contest_id = ?", userID, input.ContestID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Reminder already set for this contest"})
		return
	}

	reminder := models.Reminder{
		UserID:    userID,
		ContestID: input.ContestID,
		SendAt:    sendAt,
		Sent:      false,
	}

	if err := boot.DB.Create(&reminder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reminder"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Reminder set successfully",
		"reminder": reminder,
		"send_at":  sendAt.Format(time.RFC3339),
	})
}

//GET /reminders

func GetReminders(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))

	var reminders []models.Reminder
	if err := boot.DB.Preload("Contest").Where("user_id = ?", userID).Order("send_at ASC").Find(&reminders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reminder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":     len(reminders),
		"reminders": reminders,
	})
}

//DELETE /reminder/:id

func DeleteReminder(c *gin.Context) {
	userID := uint(c.GetFloat64("user_id"))
	reminderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reminder ID"})
		return
	}

	result := boot.DB.Where("id = ? AND user_id = ?", reminderID, userID).Delete(&models.Reminder{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reminder"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted"})
}
