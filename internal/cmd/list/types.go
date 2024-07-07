package list

// TODO
type ListScopesResponseT struct {
	Items []ScopeT `json:"items"`
}

// TODO
type ScopeT struct {
	Id      string `json:"id"`
	ScopeId string `json:"scope_id"`
	Scope   struct {
		Id            string `json:"id"`
		Type          string `json:"type"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		ParentScopeId string `json:"parent_scope_id,omitempty"`
	} `json:"scope"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// TODO
type ListTargetsResponseT struct {
	Items []TargetT `json:"items"`
}

// TODO
type TargetT struct {
	Id                     string `json:"id"`
	Name                   string `json:"name"`
	Type                   string `json:"type"`
	SessionMaxSeconds      int    `json:"session_max_seconds"`
	SessionConnectionLimit int    `json:"session_connection_limit"`
	Attributes             struct {
		DefaultPort int `json:"default_port"`
	} `json:"attributes"`
	Address string `json:"address"`
}

// AbbreviationToScopeMap represents a relationship between abbreviation and the real scope id
type AbbreviationToScopeMapT map[string]string
