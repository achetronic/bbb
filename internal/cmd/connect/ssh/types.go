package ssh

// AuthorizeSessionResponseT represents... TODO
type AuthorizeSessionResponseT struct {
	StatusCode int `json:"status_code"`
	Item       struct {
		SessionId string `json:"session_id"`
		TargetId  string `json:"target_id"`

		AuthorizationToken string `json:"authorization_token"`
		Credentials        []struct {
			Secret struct {
				Decoded struct {
					Data map[string]string `json:"data"`
				} `json:"decoded"`
			} `json:"secret"`
			Credential struct {
				PrivateKey string `json:"private_key"`
				Username   string `json:"username"`
				Password   string `json:"password"`
			} `json:"credential"`
		} `json:"credentials"`
	} `json:"item"`
}
