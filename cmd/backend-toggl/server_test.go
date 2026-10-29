package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.livingit.de/timenote"
	"go.livingit.de/timenote/remote"
)

// fakeBackend records the last call as a string.
type fakeBackend struct {
	last string
	err  error
}

func (f *fakeBackend) rec(s string) error { f.last = s; return f.err }

func (f *fakeBackend) New() error                    { return f.rec("New") }
func (f *fakeBackend) StartPrevious() error          { return f.rec("StartPrevious") }
func (f *fakeBackend) Append(l, s string) error      { return f.rec("Append " + l + " " + s) }
func (f *fakeBackend) SetDescription(l string) error { return f.rec("SetDescription " + l) }
func (f *fakeBackend) Tag(n string) error            { return f.rec("Tag " + n) }
func (f *fakeBackend) SetProjectForCurrentTimestamp(n string, a bool) error {
	if a {
		return f.rec("SetProject " + n + " auto")
	}
	return f.rec("SetProject " + n)
}
func (f *fakeBackend) Done() error { return f.rec("Done") }
func (f *fakeBackend) Current() (*timenote.TimeEntry, error) {
	return &timenote.TimeEntry{ID: 7, Note: "n"}, f.rec("Current")
}
func (f *fakeBackend) ListForDay() ([]timenote.TimeEntry, error) {
	return []timenote.TimeEntry{{ID: 1}}, f.rec("ListForDay")
}
func (f *fakeBackend) Projects() ([]timenote.Project, error) {
	return []timenote.Project{{ID: 2, Name: "p"}}, f.rec("Projects")
}
func (f *fakeBackend) CreateProject(n string) error { return f.rec("CreateProject " + n) }
func (f *fakeBackend) DeleteProject(n string) error { return f.rec("DeleteProject " + n) }
func (f *fakeBackend) Clients() ([]timenote.Client, error) {
	return []timenote.Client{{ID: 3, Name: "c"}}, f.rec("Clients")
}
func (f *fakeBackend) CreateClient(n string) error { return f.rec("CreateClient " + n) }

func TestRoundTrip(t *testing.T) {
	fake := &fakeBackend{}
	ts := httptest.NewServer(newHandler(fake, "secret"))
	defer ts.Close()
	c := remote.New(ts.URL, "secret")

	calls := []struct {
		want string
		fn   func() error
	}{
		{"New", c.New},
		{"StartPrevious", c.StartPrevious},
		{"Append a ;", func() error { return c.Append("a", ";") }},
		{"SetDescription d", func() error { return c.SetDescription("d") }},
		{"Tag t", func() error { return c.Tag("t") }},
		{"SetProject p auto", func() error { return c.SetProjectForCurrentTimestamp("p", true) }},
		{"SetProject ", func() error { return c.SetProjectForCurrentTimestamp("", false) }},
		{"Done", c.Done},
		{"CreateProject p", func() error { return c.CreateProject("p") }},
		{"DeleteProject a/b c", func() error { return c.DeleteProject("a/b c") }},
		{"CreateClient c", func() error { return c.CreateClient("c") }},
	}
	for _, call := range calls {
		if err := call.fn(); err != nil {
			t.Fatalf("%s: %v", call.want, err)
		}
		if fake.last != call.want {
			t.Errorf("got %q, want %q", fake.last, call.want)
		}
	}

	te, err := c.Current()
	if err != nil || te.ID != 7 || te.Note != "n" {
		t.Errorf("Current: %+v, %v", te, err)
	}
	entries, err := c.ListForDay()
	if err != nil || len(entries) != 1 || entries[0].ID != 1 {
		t.Errorf("ListForDay: %+v, %v", entries, err)
	}
	projects, err := c.Projects()
	if err != nil || !reflect.DeepEqual(projects, []timenote.Project{{ID: 2, Name: "p"}}) {
		t.Errorf("Projects: %+v, %v", projects, err)
	}
	clients, err := c.Clients()
	if err != nil || !reflect.DeepEqual(clients, []timenote.Client{{ID: 3, Name: "c"}}) {
		t.Errorf("Clients: %+v, %v", clients, err)
	}
}

func TestErrors(t *testing.T) {
	fake := &fakeBackend{}
	ts := httptest.NewServer(newHandler(fake, "secret"))
	defer ts.Close()

	var problem *remote.Problem

	err := remote.New(ts.URL, "wrong").New()
	if !errors.As(err, &problem) || problem.Status != http.StatusUnauthorized || problem.Code != "UNAUTHORIZED" {
		t.Errorf("wrong token: %v", err)
	}
	if fake.last != "" {
		t.Errorf("backend called without auth: %q", fake.last)
	}

	c := remote.New(ts.URL, "secret")
	err = c.CreateProject("")
	if !errors.As(err, &problem) || problem.Status != http.StatusUnprocessableEntity {
		t.Errorf("empty name: %v", err)
	}

	fake.err = fmt.Errorf("wrapped: %w", timenote.ErrNoCurrentEntry)
	err = c.Done()
	if !errors.As(err, &problem) || problem.Status != http.StatusNotFound || problem.Code != "NO_CURRENT_ENTRY" {
		t.Errorf("no current entry: %v", err)
	}

	fake.err = errors.New("toggl down")
	err = c.Done()
	if !errors.As(err, &problem) || problem.Status != http.StatusBadGateway || problem.Detail != "toggl down" {
		t.Errorf("backend error: %v", err)
	}
	if _, err := c.Current(); err == nil {
		t.Error("Current: expected error")
	}
}
