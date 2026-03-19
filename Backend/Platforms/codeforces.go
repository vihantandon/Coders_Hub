package platforms

import (
	"encoding/json"
	"net/http"

	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

type CFResponse struct {
	CFContests []models.Contest `json:"cf_future_contests"`
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
	json.NewDecoder(res.Body).Decode(&data)

	var contests []models.Contest
	for _, c := range data.CFContests {
		contests = append(contests, models.Contest{
			Name:  c.Name,
			Code:  c.Code,
			Start: c.Start,
			End:   c.End,
		})
	}

	ch <- contests
}
