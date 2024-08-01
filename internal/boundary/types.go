package boundary

// AuthorizeSessionResponseT represents... TODO
type AuthorizeSessionResponseT struct {
	StatusCode int `json:"status_code"`
	Item       struct {
		SessionId string `json:"session_id"`
		TargetId  string `json:"target_id"`
		Endpoint  string `json:"endpoint"`

		AuthorizationToken string `json:"authorization_token"`
		Credentials        []struct {
			Secret struct {
				Decoded struct {
					Data                    map[string]string `json:"data,omitempty"`
					ServiceAccountName      string            `json:"service_account_name,omitempty"`
					ServiceAccountNamespace string            `json:"service_account_namespace,omitempty"`
					ServiceAccountToken     string            `json:"service_account_token,omitempty"`
				} `json:"decoded,omitempty"`
			} `json:"secret,omitempty"`
			Credential struct {
				PrivateKey string `json:"private_key,omitempty"`
				Username   string `json:"username,omitempty"`
				Password   string `json:"password,omitempty"`
			} `json:"credential,omitempty"`
		} `json:"credentials"`
	} `json:"item"`
}

// TODO
type ConnectSessionStdoutT struct {
	Address    string `json:"address"`
	Expiration string `json:"expiration"`
	Port       int    `json:"port"`
	SessionId  string `json:"session_id"`
}
