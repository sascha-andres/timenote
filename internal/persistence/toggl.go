package persistence

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/sascha-andres/go-toggl"
	"livingit.de/code/timenote"
	"livingit.de/code/timenote/internal/cache"
)

type (
	TogglPersistor struct {
		dsn       string
		workspace int
		session   toggl.Session
		caching   *cache.Cache
	}
)

// NewToggl establishes a session to toggl api
// token constains the api token to access the toggl pai
// workspace defines the workspace to work within
func NewToggl(token string, workspace int, caching *cache.Cache) (*TogglPersistor, error) {
	res := TogglPersistor{
		dsn:       token,
		workspace: workspace,
		session:   toggl.OpenSession(token),
		caching:   caching,
	}
	toggl.DisableLog()
	return &res, res.guessWorkspace()
}

func (t *TogglPersistor) guessWorkspace() error {
	if t.workspace == 0 {
		account, err := t.session.GetAccount()
		if err != nil {
			return errors.Wrap(err, "unable to get account")
		}

		t.workspace = account.Data.Workspaces[0].ID
	}
	return nil
}

// New starts a new time entry with no description
func (t *TogglPersistor) New() error {
	_, err := t.session.StartTimeEntry("")
	if err != nil {
		return errors.Wrap(err, "Unable to start a new entry")
	}
	return nil
}

// Append adds text to the description separated by ;
func (t *TogglPersistor) Append(line, separator string) error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "Unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "unable to get running time entry from toggl")
	}
	if te.Description == "" {
		te.Description = line
	} else {
		te.Description = fmt.Sprintf("%s%s%s", te.Description, separator, line)
	}
	_, err = t.session.UpdateTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "unable to update time entry in toggl")
	}
	return nil
}

// Tag toggle the tags associated with the running time entry
// name is the name of the tag
func (t *TogglPersistor) Tag(name string) error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "unable to get running time entry from toggl")
	}
	if te.HasTag(name) {
		te.RemoveTag(name)
	} else {
		te.AddTag(name)
	}
	_, err = t.session.UpdateTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "unable to update time entry in toggl")
	}
	return nil
}

// Done ends the currently running time entry
func (t *TogglPersistor) Done() error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "unable to get running time entry from toggl")
	}
	_, err = t.session.StopTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "unable to stop running time entry in toggl")
	}
	return nil
}

// Current returns the currently running time entry
func (t *TogglPersistor) Current() (*timenote.TimeEntry, error) {
	account, err := t.session.GetAccount()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get running time entry from toggl")
	}
	var res timenote.TimeEntry
	res.Note = te.Description
	res.ID = te.ID
	res.Tag = fmt.Sprintf("%v", te.Tags)
	res.Start = *te.Start
	res.Duration = te.Duration
	if te.Pid > 0 {
		p, err := t.GetProject(te.Pid)
		if err != nil {
			res.Project = "- unknown project -"
		} else {
			res.Project = p.Name
			if p.Cid > 0 {
				c, err := t.GetClientByID(p.Cid)
				if err != nil {
					res.Client = "- unknown client -"
				} else {
					res.Client = c.Name
				}
			}
		}
	}
	return &res, nil
}

// GetClientByID gets all clients and returns the one with the given ID
func (t *TogglPersistor) GetClientByID(clientID int) (*timenote.Client, error) {
	cs, err := t.Clients()
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		if c.ID == clientID {
			return &c, nil
		}
	}
	return nil, nil
}

// GetProject returns the given project
func (t *TogglPersistor) GetProject(projectID int) (*toggl.Project, error) {
	p, err := t.session.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getCurrentTimeEntry(account toggl.Account) (*toggl.TimeEntry, error) {
	for _, te := range account.Data.TimeEntries {
		if nil == te.Stop {
			return &te, nil
		}
	}
	return nil, fmt.Errorf("no current time entry")
}

// SetProjectForCurrentTimestamp apply project to running time entry
func (t *TogglPersistor) SetProjectForCurrentTimestamp(name string) error {
	var (
		account   toggl.Account
		projectID int
		err       error
	)

	account, err = t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "unable to get account")
	}
	if name == "" {
		projectID = 0
	} else {
		projectID, err = t.getProjectID(name)
		if projectID == 0 {
			projectID, err = t.createProject(account, name)
			if err != nil {
				return errors.Wrap(err, "unable to create project")
			}
		}
	}

	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "unable to get running time entry from toggl")
	}
	te.Pid = projectID
	_, err = t.session.UpdateTimeEntry(*te)
	return err
}

func (t *TogglPersistor) createProject(account toggl.Account, name string) (int, error) {
	res, err := t.session.CreateProject(name, t.workspace)
	if err != nil {
		return 0, err
	}
	return res.ID, nil
}

// CreateProject creates a new project within the workspace
func (t *TogglPersistor) CreateProject(name string) error {
	id, err := t.getProjectID(name)
	if err != nil {
		return err
	}
	if id == 0 {
		account, err := t.session.GetAccount()
		if err != nil {
			return err
		}
		_, err = t.createProject(account, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteProject removes a project from the workspace
func (t *TogglPersistor) DeleteProject(name string) error {
	var (
		project *toggl.Project
		id      int
		err     error
	)

	if id, err = strconv.Atoi(name); !(err == nil && id != 0) {
		id, err = t.getProjectID(name)
	}

	if err != nil {
		return err
	}

	project, err = t.session.GetProject(id)

	if nil == project {
		return errors.New("no such project")
	}

	_, err = t.session.DeleteProject(toggl.Project{
		Wid: project.Wid,
		ID:  project.ID,
		Cid: project.Cid,
	})
	return err
}

func (t *TogglPersistor) getProjectID(name string) (int, error) {
	account, err := t.session.GetAccount()
	if err != nil {
		return 0, errors.Wrap(err, "unable to get account")
	}

	for _, prj := range account.Data.Projects {
		if prj.Name == name {
			return prj.ID, nil
		}
	}

	return 0, nil
}

// Clients ereturn all clients
func (t *TogglPersistor) Clients() ([]timenote.Client, error) {
	clients, err := t.session.GetClients()
	if err != nil {
		return nil, err
	}
	var result = make([]timenote.Client, 0)
	for _, c := range clients {
		result = append(result, timenote.Client{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Notes,
		})
	}
	return result, nil
}

// NewClient creates a new client
func (t *TogglPersistor) NewClient(name string) error {
	_, err := t.session.CreateClient(name, t.workspace)
	return err
}

// ListForDay returns all time entries for a day
func (t *TogglPersistor) ListForDay() ([]timenote.TimeEntry, error) {
	year, month, day := time.Now().Date()
	loc, _ := time.LoadLocation("")
	startDate := time.Date(year, month, day, 0, 0, 0, 0, loc)
	endDate := time.Date(year, month, day, 23, 59, 59, 0, loc)
	entries, err := t.session.GetTimeEntries(startDate, endDate)
	if err != nil {
		return nil, err
	}
	result := make([]timenote.TimeEntry, 0)
	for _, entry := range entries {
		te := timenote.TimeEntry{
			ID:       entry.ID,
			Tag:      fmt.Sprintf("%v", entry.Tags),
			Note:     entry.Description,
			Start:    *entry.Start,
			Stop:     entry.Stop,
			Duration: entry.Duration,
		}
		if entry.Pid > 0 {
			p, err := t.GetProject(entry.Pid)
			if err != nil {
				te.Project = "- unknown project -"
			} else {
				te.Project = p.Name
				if p.Cid > 0 {
					c, err := t.GetClientByID(p.Cid)
					if err != nil {
						te.Client = "- unknown client -"
					} else {
						te.Client = c.Name
					}
				}
			}
		}
		result = append(result, te)
	}
	return result, nil
}

// Projects returns a list of all projects
func (t *TogglPersistor) Projects() ([]timenote.Project, error) {
	projects, err := t.session.GetProjects(t.workspace)
	if err != nil {
		return nil, err
	}
	result := make([]timenote.Project, 0)
	for _, prj := range projects {
		result = append(result, timenote.Project{
			ID:          prj.ID,
			WorkspaceID: prj.Wid,
			ClientID:    prj.Cid,
			Name:        prj.Name,
			Billable:    prj.Billable,
		})
	}
	return result, nil
}
