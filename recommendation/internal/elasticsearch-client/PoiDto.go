package elasticsearch_client

/*
*
ES POI Document
*/
type PoiDocument struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	Coordinates []float64 `json:"coordinates"`
	Address     string    `json:"address"`
}

type GeoType struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
