package metadata

type ImageResizeRequest struct {
	Url    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}

type ImageToWebPRequest struct {
	Url string `json:"url"`
}

type ImageDrawRequest struct {
	Url string `json:"url"`
	X0  int64  `json:"x0"`
	Y0  int64  `json:"y0"`
	X1  int64  `json:"x1"`
	Y1  int64  `json:"y1"`
}
