package auth

// ResponseT represents TODO
type ResponseT struct {
	Item struct {
		Attributes struct {
			Token string `json:"token"`
		} `json:"attributes"`
	} `json:"item"`
}
