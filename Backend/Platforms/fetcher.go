package platforms

import (
	"sync"

	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

func FetchAndStore(logger *zap.SugaredLogger) {
	type platformFetcher struct {
		name string
		fn   func(*zap.SugaredLogger, chan []models.Contest)
	}

	fetchers := []platformFetcher{
		{"CodeChef", FetchCodeChef},
		{"CodeForces", FetchCodeForces},
		{"Leetcode", FetchLeetcode},
	}

	ch := make(chan []models.Contest, len(fetchers))
	var wg sync.WaitGroup

	for _, f := range fetchers {
		wg.Add(1)
		go func(pf platformFetcher) {
			defer wg.Done()
			pf.fn(logger, ch)
		}(f)
	}

	wg.Wait()
	close(ch)

	var total int
	for contests := range ch {
		for _, c := range contests {

			result := boot.DB.Where(models.Contest{Code: c.Code, Platform: c.Platform}).Assign(models.Contest{
				Name:  c.Name,
				Start: c.Start,
				End:   c.End,
			}).FirstOrCreate(&c)

			if result.Error != nil {
				logger.Errorf("Failed to upsert contests %s: %v", c.Name, result.Error)
			} else {
				total++
			}
		}
	}

	logger.Info("Fetched and stored %d contests", total)
}
