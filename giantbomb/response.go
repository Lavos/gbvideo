package giantbomb

type Response struct {
	StatusCode int64 `json:"status_code"`
	Error string `json:"error"`
	TotalResults int64 `json:"number_of_total_results"`
	PageResults int64 `json:"number_of_page_results"`
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
}
