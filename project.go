package timenote

type (
	Project struct {
		ID          int `json:"id"`
		WorkspaceID int `json:"workspace_id"`
		// ClientID is 0 for a project without a client
		ClientID  int    `json:"client_id,omitempty"`
		Name      string `json:"name"`
		Billable  bool   `json:"billable"`
		IsPrivate bool   `json:"is_private"`
	}
)
