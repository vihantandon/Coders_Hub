package platforms

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

type CCContest struct {
	ContestCode  string `json:"contest_code"`
	ContestName  string `json:"contest_name"`
	ContestStart string `json:"contest_start_date"`
	ContestEnd   string `json:"contest_end_date"`
}
type CCResponse struct {
	FutureContests []CCContest `json:"future_contests"`
}

func FetchCodeChef(logger *zap.SugaredLogger, ch chan []models.Contest) {
	url := "https://www.codechef.com/api/list/contests/all"

	res, err := http.Get(url)
	if err != nil {
		logger.Errorf("Error fetching data from codechef API: %v", err)
		ch <- nil
		return
	}

	defer res.Body.Close()

	var data CCResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		logger.Errorf("Error decoding codechef response %v", err)
		ch <- nil
		return
	}

	var contests []models.Contest

	for _, c := range data.FutureContests {
		layout := "02 Jan 2006  15:04:05"

		startTime, err1 := time.Parse(layout, c.ContestStart)
		endTime, err2 := time.Parse(layout, c.ContestEnd)

		if err1 != nil || err2 != nil {
			logger.Errorf("Error parsing time: %v %v", err1, err2)
			continue
		}

		formattedLayout := "2006-01-02 15:04:05"
		contests = append(contests, models.Contest{
			Name:     c.ContestName,
			Code:     c.ContestCode,
			Platform: "CodeChef",
			Start:    startTime.Format(formattedLayout),
			End:      endTime.Format(formattedLayout),
		})
	}

	ch <- contests
}
