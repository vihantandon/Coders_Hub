package platforms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/vihantandon/Coders_Hub/models"
	"go.uber.org/zap"
)

type LCContest struct {
	Title     string `json:"title"`
	TitleSlug string `json:"titleSlug"`
	StartTime int64  `json:"startTime"`
	Duration  int    `json:"duration"`
}

type LCData struct {
	TopTowContests []LCContest `json:"topTwoContests"`
	AllContests    []LCContest `json:"allContests"`
}

type LCResponse struct {
	Data LCData `json:"data"`
}

func FetchLeetcode(logger *zap.SugaredLogger, ch chan []models.Contest) {
	query := `{
		"query": "{allContests {title titleSlug startTime duration}}"
	}`

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql",
		bytes.NewBuffer([]byte(query)))

	if err != nil {
		logger.Errorf("Error creating leetcode request: %v", err)
		ch <- nil
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Error fetching Leetcode: %v", err)
		ch <- nil
		return
	}

	defer res.Body.Close()

	var data LCResponse

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		logger.Errorf("Error decoding Leetcode response: %v", err)
		ch <- nil
		return
	}

	now := time.Now().Unix()

	var contests []models.Contest
	for _, c := range data.Data.AllContests {
		if c.StartTime < now {
			continue
		}

		start := time.Unix(c.StartTime, 0).UTC()
		end := time.Unix(c.StartTime+int64(c.Duration), 0).UTC()
		contests = append(contests, models.Contest{
			Name:     c.Title,
			Code:     c.TitleSlug,
			Platform: "Leetcode",
			Start:    start.Format("2006-01-02 15:04:05"),
			End:      end.Format("2006-01-02 15:04:05"),
		})
	}

	ch <- contests
}
