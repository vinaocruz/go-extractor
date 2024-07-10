package dto

type Negociation struct {
	Ticker           string  `json:"ticker"`
	Max_range_value  float32 `json:"max_range_value"`
	Max_daily_volume int     `json:"max_daily_volume"`
}
