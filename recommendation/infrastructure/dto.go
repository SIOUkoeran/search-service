package infrastructure

type ResponseEmbedQuery struct {
	Vector []float64 `json:"vector"`
}

type RequestEmbedQuery struct {
	Query string `json:"query"`
}
