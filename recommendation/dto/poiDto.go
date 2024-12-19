package dto

const (
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
)

type (
	TitleSearchRequest struct {
		Title string `json:"title" form:"title"`
	}

	PoiEntity struct {
		Title    string   `json:"title" form:"title"`
		Address  string   `json:"address" form:"address"`
		Location Location `json:"location" form:"location"`
	}

	Location struct {
		Lat float32 `json:"lat" form:"lat"`
		Lon float32 `json:"lon" form:"lon"`
	}
)
