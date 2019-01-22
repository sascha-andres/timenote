package toggldriver

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sascha-andres/go-toggl"
	"livingit.de/code/timenote"
	"livingit.de/code/timenote/persistence"
)

type (
	TogglPersistor struct {
		dsn     string
		session toggl.Session
	}
)

func NewToggl(dsn string) (persistence.Persistor, error) {
	var res TogglPersistor
	res.dsn = dsn
	res.session = toggl.OpenSession(dsn)
	toggl.DisableLog()
	return &res, nil
}

func (t *TogglPersistor) New() error {
	_, err := t.session.StartTimeEntry("")
	if err != nil {
		return errors.Wrap(err, "Unable to start a new entry")
	}
	return nil
}

func (t *TogglPersistor) Append(line string) error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "Unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "Unable to get running time entry from toggl")
	}
	if te.Description == "" {
		te.Description = line
	} else {
		te.Description = fmt.Sprintf("%s;%s", te.Description, line)
	}
	_, err = t.session.UpdateTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "Unable to update time entry in toggl")
	}
	return nil
}

func (t *TogglPersistor) Tag(name string) error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "Unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "Unable to get running time entry from toggl")
	}
	if te.HasTag(name) {
		te.RemoveTag(name)
	} else {
		te.AddTag(name)
	}
	_, err = t.session.UpdateTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "Unable to update time entry in toggl")
	}
	return nil
}

func (t *TogglPersistor) Done() error {
	account, err := t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "Unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "Unable to get running time entry from toggl")
	}
	_, err = t.session.StopTimeEntry(*te)
	if err != nil {
		return errors.Wrap(err, "Unable to stop running time entry in toggl")
	}
	return nil
}

func (t *TogglPersistor) Close() error {
	return nil
}

/*func (t *TogglPersistor) ListForDay(delta int) ([]timenote.TimeEntry, error) {
	return nil, errors.New("Not yet implemented")
}*/

func (t *TogglPersistor) Current() (*timenote.TimeEntry, error) {
	account, err := t.session.GetAccount()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get account")
	}
	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get running time entry from toggl")
	}
	var res timenote.TimeEntry
	res.Note = te.Description
	res.ID = te.ID
	res.Tag = fmt.Sprintf("%v", te.Tags)
	res.Start = *te.Start
	res.Duration = te.Duration
	return &res, nil
}

func getCurrentTimeEntry(account toggl.Account) (*toggl.TimeEntry, error) {
	for _, te := range account.Data.TimeEntries {
		if nil == te.Stop {
			return &te, nil
		}
	}
	return nil, fmt.Errorf("No current time entry")
}

func (t *TogglPersistor) Project(name string) error {
	var (
		account   toggl.Account
		projectID int
		err       error
	)

	account, err = t.session.GetAccount()
	if err != nil {
		return errors.Wrap(err, "Unable to get account")
	}
	if name == "" {
		projectID = 0
	} else {
		projectID, err = t.getProjectID(name)
		if projectID == 0 {
			projectID, err = t.createProject(account, name)
			if err != nil {
				return errors.Wrap(err, "Unable to create project")
			}
		}
	}

	te, err := getCurrentTimeEntry(account)
	if err != nil {
		return errors.Wrap(err, "Unable to get running time entry from toggl")
	}
	te.Pid = projectID
	_, err = t.session.UpdateTimeEntry(*te)
	return err
}

func (t *TogglPersistor) createProject(account toggl.Account, name string) (int, error) {
	res, err := t.session.CreateProject(name, account.Data.Workspaces[0].ID)
	if err != nil {
		return 0, err
	}
	return res.ID, nil
}

func (t *TogglPersistor) getProjectID(name string) (int, error) {
	account, err := t.session.GetAccount()
	if err != nil {
		return 0, errors.Wrap(err, "Unable to get account")
	}

	for _, prj := range account.Data.Projects {
		if prj.Name == name {
			return prj.ID, nil
		}
	}

	return 0, nil
}

func (t *TogglPersistor) GetWebsite() (bool, string, error) {
	return true, "https://toggl.com/app/timer", nil
}

func (t *TogglPersistor) Client() ([]timenote.Client, error) {
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

func (t *TogglPersistor) NewClient(name string) error {
	account, err := t.session.GetAccount()
	if err != nil {
		return err
	}
	_, err = t.session.CreateClient(name, account.Data.Workspaces[0].ID)
	return err
}
