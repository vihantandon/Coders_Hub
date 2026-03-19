package platforms

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/vihantandon/Coders_Hub/boot"
)

type Contest struct {
	Code  string `json:"contest_code"`
	Name  string `json:"contest_name"`
	Start string `json:"contest_start_date"`
	End   string `json:"contest_end_date"`
}

type Response struct {
	FutureContests []Contest `json:"future_contests"`
}

func main() {
	logger := boot.InitializeApp()
	url := "https://www.codechef.com/api/list/contests/all"

	res, err := http.Get(url)
	if err != nil {
		logger.Errorf("Error fetching data from codechef API: %v", err)
	}

	defer res.Body.Close()

	var data Response
	json.NewDecoder(res.Body).Decode(&data)

	var wg sync.WaitGroup

	for _, c := range data.FutureContests {

		wg.Add(1)
		go func(contest Contest) {

			defer wg.Done()
			logger.Infof(
				"Name: %s | Code: %s | Start: %s | End: %s",
				contest.Name,
				contest.Code,
				contest.Start,
				contest.End,
			)
		}(c) // passes c into goroutine
	}

	wg.Wait()
}
