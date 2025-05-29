package data

type SavedResponse struct {
	Method   string `json:"method"`
	Url      string `json:"url"`
	Status   int    `json:"status"`
	Response any    `json:"response"`
}
