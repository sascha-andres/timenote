// Package remote is a timenote.Backend that calls a backend-toggl server.
package remote

import "fmt"

// Request bodies of the backend-toggl /v1 API.
type (
	NewTimeEntryRequest struct {
		// FromPrevious starts the entry based on the last finished time entry
		FromPrevious bool `json:"from_previous"`
	}

	SetDescriptionRequest struct {
		Description string `json:"description"`
	}

	AppendDescriptionRequest struct {
		Line      string `json:"line"`
		Separator string `json:"separator"`
	}

	// UpdateCurrentRequest changes the running time entry
	UpdateCurrentRequest struct {
		// Stop stops the running time entry
		Stop bool `json:"stop"`
	}

	ToggleTagRequest struct {
		Toggle string `json:"toggle"`
	}

	SetProjectRequest struct {
		// Name is the project name; empty removes the project
		Name              string `json:"name"`
		AutoCreateProject bool   `json:"auto_create_project"`
	}

	NameRequest struct {
		Name string `json:"name"`
	}

	// Problem is an RFC 9457 error response.
	Problem struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Status int    `json:"status"`
		Detail string `json:"detail,omitempty"`
		Code   string `json:"code"`
	}
)

// ProblemContentType is the media type of a Problem response.
const ProblemContentType = "application/problem+json"

func (p *Problem) Error() string {
	return fmt.Sprintf("%s (%d %s): %s", p.Title, p.Status, p.Code, p.Detail)
}
