package main

import (
	"sync"

	platforms "github.com/vihantandon/Coders_Hub/Platforms"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
)

func main() {
	logger := boot.InitializeApp()

	var wg sync.WaitGroup
	ch := make(chan []models.Contest, 2)

	wg.Add(2)

	go func() {
		defer wg.Done()
		platforms.FetchCodeChef(logger, ch)
	}()

	go func() {
		defer wg.Done()
		platforms.FetchCodeForces(logger, ch)
	}()

	wg.Wait()
	close(ch)

	for contests := range ch {
		for _, c := range contests {
			logger.Infof("Contests: %s", c.Name)
		}
	}
}
