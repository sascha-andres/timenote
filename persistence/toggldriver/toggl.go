package toggldriver

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/jason0x43/go-toggl"
	"github.com/sascha-andres/timenote"
	"github.com/sascha-andres/timenote/persistence"
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
	t.session.UpdateTimeEntry(*te)
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

func (t *TogglPersistor) ListForDay(delta int) ([]timenote.TimeEntry, error) {
	return nil, errors.New("Not yet implemented")
}

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
	return nil, nil
}

func getCurrentTimeEntry(account toggl.Account) (*toggl.TimeEntry, error) {
	for _, te := range account.Data.TimeEntries {
		if nil == te.Stop {
			return &te, nil
		}
	}
	return nil, fmt.Errorf("No current time entry")
}
