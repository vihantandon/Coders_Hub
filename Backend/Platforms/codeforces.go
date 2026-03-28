package platforms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

type CFContest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
	DurationSeconds  int64  `json:"durationSeconds"`
}

type CFResponse struct {
	Status string      `json:"status"`
	Result []CFContest `json:"result"`
}

func FetchCodeForces(logger *zap.SugaredLogger, ch chan []models.Contest) {
	url := "https://codeforces.com/api/contest.list"

	res, err := http.Get(url)
	if err != nil {
		logger.Errorf("Error fetching data from codeforces API: %v", err)
		ch <- nil
		return
	}

	defer res.Body.Close()

	var data CFResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		logger.Errorf("Error decoding codeforces response: %v", err)
		ch <- nil
		return
	}

	var contests []models.Contest
	now := time.Now().Unix()

	for _, c := range data.Result {
		if c.StartTimeSeconds < now {
			continue
		}

		starTime := time.Unix(c.StartTimeSeconds, 0).UTC()
		endTime := time.Unix(c.StartTimeSeconds+c.DurationSeconds, 0).UTC()
		contests = append(contests, models.Contest{
			Name:     c.Name,
			Code:     fmt.Sprintf("%d", c.ID),
			Platform: "Codeforces",
			Start:    starTime.Format("2006-01-02 15:04:05"),
			End:      endTime.Format("2006-01-02 15:04:05"),
		})
	}

	ch <- contests
}
