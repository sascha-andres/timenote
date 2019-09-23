package timenote

type (
	Project struct {
		ID          int
		WorkspaceID int
		ClientID    int
		Name        string
		Billable    bool
		IsPrivate   bool
	}
)
