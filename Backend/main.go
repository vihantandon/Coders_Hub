package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vihantandon/Coders_Hub/boot"
	"github.com/vihantandon/Coders_Hub/routes"
	"github.com/vihantandon/Coders_Hub/scheduler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	// var wg sync.WaitGroup
	// ch := make(chan []models.Contest, 2)

	// wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	platforms.FetchCodeChef(logger, ch)
	// }()

	// go func() {
	// 	defer wg.Done()
	// 	platforms.FetchCodeForces(logger, ch)
	// }()

	// wg.Wait()
	// close(ch)

	// for contests := range ch {
	// 	for _, c := range contests {
	// 		logger.Infof("Contests: %s", c.Name)
	// 	}
	// }

	logger := boot.InitializeApp()
	boot.InitDB()

	scheduler.StartContestScheduler(logger)

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
