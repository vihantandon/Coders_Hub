package models

type Contest struct {
	Code  string `json:"contest_code"`
	Name  string `json:"contest_name"`
	Start string `json:"contest_start_date"`
	End   string `json:"contest_end_date"`
}
