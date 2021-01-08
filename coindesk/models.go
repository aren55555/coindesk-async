package coindesk

import "time"

// Price models the repeated pricing structure of the API response.
type Price struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Description string  `json:"description"`
	Rate        float64 `json:"rate_float"`
}

// Response models the outermost JSON structure of the API response.
type Response struct {
	Time struct {
		Updated string `json:"updatedISO"`
	} `json:"time"`
	BPI struct {
		USD Price `json:"USD"`
		GBP Price `json:"GBP"`
		EUR Price `json:"EUR"`
	} `json:"bpi"`
}

// Value records the time at which the Go program fetched the result.
type Value struct {
	FetchTime time.Time
	Result    Response
}
