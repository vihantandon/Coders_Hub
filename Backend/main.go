package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	platforms "github.com/vihantandon/Coders_Hub/Platforms"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/models"
	"github.com/vihantandon/Coders_Hub/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

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

	boot.InitDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
