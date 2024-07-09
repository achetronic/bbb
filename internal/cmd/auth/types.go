package auth

// ResponseT represents TODO
type ResponseT struct {
	StatusCode int `json:"status_code"`
	Item       struct {
		Attributes struct {
			Token string `json:"token"`
		} `json:"attributes"`
	} `json:"item"`
}
