package scheduler

import (
	"time"

	platforms "github.com/vihantandon/Coders_Hub/Platforms"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/mailer"
	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

func StartContestScheduler(logger *zap.SugaredLogger) {
	logger.Info("Running initial contest fetch")
	platforms.FetchAndStore(logger)

	//fetch after 1 hour
	ticker := time.NewTicker(5 * time.Hour)
	go func() {
		for range ticker.C {
			logger.Info("Scheduler: fetching contests...")
			platforms.FetchAndStore(logger)
		}
	}()

	reminderTicker := time.NewTicker(1 * time.Minute)
	go func() {
		for range reminderTicker.C {
			sendPendingReminder(logger)
		}
	}()
}

func sendPendingReminder(logger *zap.SugaredLogger) {
	now := time.Now().UTC()

	var reminders []models.Reminder
	err := boot.DB.Preload("User").
		Preload("Contest").
		Where("sent = false AND send_at <= ?", now).
		Find(&reminders).Error

	if err != nil {
		logger.Errorf("failed to query pending reminders: %v", err)
		return
	}

	if len(reminders) == 0 {
		return
	}

	logger.Infof("Sending %d pending reminder(s)...", len(reminders))

	for _, r := range reminders {
		payload := mailer.EmailPayload{
			To:           r.User.Email,
			UserName:     r.User.Name,
			ContestName:  r.Contest.Name,
			Platform:     r.Contest.Platform,
			ContestStart: r.Contest.Start,
		}

		if err := mailer.SendReminderEmail(payload); err != nil {
			logger.Errorf("Failed to send reminder %d to %s: %v", r.ID, r.User.Email, err)
			continue
		}

		if err := boot.DB.Model(&r).Update("sent", true).Error; err != nil {
			logger.Error("failed to mark reminder %d as sent: %v", r.ID, err)
		} else {
			logger.Infof("Reminder sent to %s for contest: %s", r.User.Email, r.Contest.Name)
		}
	}
}
