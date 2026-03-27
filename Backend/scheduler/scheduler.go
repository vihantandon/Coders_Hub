package scheduler

import (
	"time"

	platforms "github.com/vihantandon/Coders_Hub/Platforms"
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
}
