package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.livingit.de/timenote"
)

// Client implements timenote.Backend through the backend-toggl REST API.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

var _ timenote.Backend = (*Client)(nil)

// New returns a client for the server at baseURL (for example
// http://127.0.0.1:8141). token is sent in the Authorization header.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// do sends body as JSON (if not nil) and decodes the response into out (if not nil).
func (c *Client) do(method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		problem := &Problem{Status: resp.StatusCode}
		if err := json.NewDecoder(resp.Body).Decode(problem); err != nil {
			return fmt.Errorf("%s %s: %s", method, path, resp.Status)
		}
		return problem
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) New() error {
	return c.do(http.MethodPost, "/v1/timeentries", NewTimeEntryRequest{}, nil)
}

func (c *Client) StartPrevious() error {
	return c.do(http.MethodPost, "/v1/timeentries", NewTimeEntryRequest{FromPrevious: true}, nil)
}

func (c *Client) Append(line, separator string) error {
	return c.do(http.MethodPatch, "/v1/timeentries/current/description", AppendDescriptionRequest{Line: line, Separator: separator}, nil)
}

func (c *Client) SetDescription(line string) error {
	return c.do(http.MethodPut, "/v1/timeentries/current/description", SetDescriptionRequest{Description: line}, nil)
}

func (c *Client) Tag(name string) error {
	return c.do(http.MethodPatch, "/v1/timeentries/current/tags", ToggleTagRequest{Toggle: name}, nil)
}

func (c *Client) SetProjectForCurrentTimestamp(name string, autoCreateProject bool) error {
	return c.do(http.MethodPut, "/v1/timeentries/current/project", SetProjectRequest{Name: name, AutoCreateProject: autoCreateProject}, nil)
}

func (c *Client) Done() error {
	return c.do(http.MethodPatch, "/v1/timeentries/current", UpdateCurrentRequest{Stop: true}, nil)
}

func (c *Client) Current() (*timenote.TimeEntry, error) {
	var te timenote.TimeEntry
	if err := c.do(http.MethodGet, "/v1/timeentries/current", nil, &te); err != nil {
		return nil, err
	}
	return &te, nil
}

func (c *Client) ListForDay() ([]timenote.TimeEntry, error) {
	var entries []timenote.TimeEntry
	if err := c.do(http.MethodGet, "/v1/timeentries", nil, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) Projects() ([]timenote.Project, error) {
	var projects []timenote.Project
	if err := c.do(http.MethodGet, "/v1/projects", nil, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *Client) CreateProject(name string) error {
	return c.do(http.MethodPost, "/v1/projects", NameRequest{Name: name}, nil)
}

func (c *Client) DeleteProject(nameOrID string) error {
	return c.do(http.MethodDelete, "/v1/projects/"+url.PathEscape(nameOrID), nil, nil)
}

func (c *Client) Clients() ([]timenote.Client, error) {
	var clients []timenote.Client
	if err := c.do(http.MethodGet, "/v1/clients", nil, &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

func (c *Client) CreateClient(name string) error {
	return c.do(http.MethodPost, "/v1/clients", NameRequest{Name: name}, nil)
}
