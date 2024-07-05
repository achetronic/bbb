package list

import "time"

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

// AbbreviationToScopeMap represents a relationship between abbreviation and the real scope id
type AbbreviationToScopeMapT map[string]string

// TODO
type T struct {
	StatusCode int `json:"status_code"`
	Items      []struct {
		Id      string `json:"id"`
		ScopeId string `json:"scope_id"`
		Scope   struct {
			Id            string `json:"id"`
			Type          string `json:"type"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			ParentScopeId string `json:"parent_scope_id"`
		} `json:"scope"`
		Name                   string    `json:"name"`
		CreatedTime            time.Time `json:"created_time"`
		UpdatedTime            time.Time `json:"updated_time"`
		Version                int       `json:"version"`
		Type                   string    `json:"type"`
		SessionMaxSeconds      int       `json:"session_max_seconds"`
		SessionConnectionLimit int       `json:"session_connection_limit"`
		Attributes             struct {
			DefaultPort int `json:"default_port"`
		} `json:"attributes"`
		AuthorizedActions []string `json:"authorized_actions"`
		Address           string   `json:"address"`
	} `json:"items"`
}
