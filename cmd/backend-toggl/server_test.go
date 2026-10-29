package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"go.livingit.de/timenote"
	"go.livingit.de/timenote/remote"
)

// fakeBackend records the last call as a string and counts the calls.
type fakeBackend struct {
	last  string
	calls int
	err   error
}

func (f *fakeBackend) rec(s string) error { f.last = s; f.calls++; return f.err }

func (f *fakeBackend) New(context.Context) error           { return f.rec("New") }
func (f *fakeBackend) StartPrevious(context.Context) error { return f.rec("StartPrevious") }
func (f *fakeBackend) Append(_ context.Context, l, s string) error {
	return f.rec("Append " + l + " " + s)
}
func (f *fakeBackend) SetDescription(_ context.Context, l string) error {
	return f.rec("SetDescription " + l)
}
func (f *fakeBackend) Tag(_ context.Context, n string) error { return f.rec("Tag " + n) }
func (f *fakeBackend) SetProjectForCurrentTimestamp(_ context.Context, n string, a bool) error {
	if a {
		return f.rec("SetProject " + n + " auto")
	}
	return f.rec("SetProject " + n)
}
func (f *fakeBackend) Done(context.Context) error { return f.rec("Done") }
func (f *fakeBackend) Current(context.Context) (*timenote.TimeEntry, error) {
	return &timenote.TimeEntry{ID: 7, Note: "n"}, f.rec("Current")
}
func (f *fakeBackend) ListForDay(context.Context) ([]timenote.TimeEntry, error) {
	return []timenote.TimeEntry{{ID: 1}}, f.rec("ListForDay")
}
func (f *fakeBackend) Projects(context.Context) ([]timenote.Project, error) {
	return []timenote.Project{{ID: 2, Name: "p"}}, f.rec("Projects")
}
func (f *fakeBackend) CreateProject(_ context.Context, n string) error {
	return f.rec("CreateProject " + n)
}
func (f *fakeBackend) DeleteProject(_ context.Context, n string) error {
	return f.rec("DeleteProject " + n)
}
func (f *fakeBackend) Clients(context.Context) ([]timenote.Client, error) {
	return []timenote.Client{{ID: 3, Name: "c"}}, f.rec("Clients")
}
func (f *fakeBackend) CreateClient(_ context.Context, n string) error {
	return f.rec("CreateClient " + n)
}
func (f *fakeBackend) UpdateCache(context.Context) error { return f.rec("UpdateCache") }
func (f *fakeBackend) CacheStatus(context.Context) ([]timenote.CacheStatus, error) {
	return []timenote.CacheStatus{{Name: "Project"}}, f.rec("CacheStatus")
}

func TestRoundTrip(t *testing.T) {
	ctx := t.Context()
	fake := &fakeBackend{}
	ts := httptest.NewServer(newHandler(fake, "secret", 0))
	defer ts.Close()
	c := remote.New(ts.URL, "secret")

	calls := []struct {
		want string
		fn   func() error
	}{
		{"New", func() error { return c.New(ctx) }},
		{"StartPrevious", func() error { return c.StartPrevious(ctx) }},
		{"Append a ;", func() error { return c.Append(ctx, "a", ";") }},
		{"SetDescription d", func() error { return c.SetDescription(ctx, "d") }},
		{"Tag t", func() error { return c.Tag(ctx, "t") }},
		{"SetProject p auto", func() error { return c.SetProjectForCurrentTimestamp(ctx, "p", true) }},
		{"SetProject ", func() error { return c.SetProjectForCurrentTimestamp(ctx, "", false) }},
		{"Done", func() error { return c.Done(ctx) }},
		{"CreateProject p", func() error { return c.CreateProject(ctx, "p") }},
		{"DeleteProject a/b c", func() error { return c.DeleteProject(ctx, "a/b c") }},
		{"CreateClient c", func() error { return c.CreateClient(ctx, "c") }},
		{"UpdateCache", func() error { return c.UpdateCache(ctx) }},
	}
	for _, call := range calls {
		if err := call.fn(); err != nil {
			t.Fatalf("%s: %v", call.want, err)
		}
		if fake.last != call.want {
			t.Errorf("got %q, want %q", fake.last, call.want)
		}
	}

	te, err := c.Current(ctx)
	if err != nil || te.ID != 7 || te.Note != "n" {
		t.Errorf("Current: %+v, %v", te, err)
	}
	entries, err := c.ListForDay(ctx)
	if err != nil || len(entries) != 1 || entries[0].ID != 1 {
		t.Errorf("ListForDay: %+v, %v", entries, err)
	}
	projects, err := c.Projects(ctx)
	if err != nil || !reflect.DeepEqual(projects, []timenote.Project{{ID: 2, Name: "p"}}) {
		t.Errorf("Projects: %+v, %v", projects, err)
	}
	clients, err := c.Clients(ctx)
	if err != nil || !reflect.DeepEqual(clients, []timenote.Client{{ID: 3, Name: "c"}}) {
		t.Errorf("Clients: %+v, %v", clients, err)
	}
	status, err := c.CacheStatus(ctx)
	if err != nil || !reflect.DeepEqual(status, []timenote.CacheStatus{{Name: "Project"}}) {
		t.Errorf("CacheStatus: %+v, %v", status, err)
	}
}

func TestReadCache(t *testing.T) {
	ctx := t.Context()
	fake := &fakeBackend{}
	ts := httptest.NewServer(newHandler(fake, "secret", time.Hour))
	defer ts.Close()
	c := remote.New(ts.URL, "secret")

	steps := []struct {
		name  string
		fn    func() error
		calls int
	}{
		{"first read", func() error { _, err := c.ListForDay(ctx); return err }, 1},
		{"cached read", func() error { _, err := c.ListForDay(ctx); return err }, 1},
		{"other key", func() error { _, err := c.Current(ctx); return err }, 2},
		{"write clears", func() error { return c.Done(ctx) }, 3},
		{"read after write", func() error { _, err := c.ListForDay(ctx); return err }, 4},
	}
	for _, step := range steps {
		if err := step.fn(); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		if fake.calls != step.calls {
			t.Errorf("%s: %d backend calls, want %d", step.name, fake.calls, step.calls)
		}
	}

	// errors other than "no current entry" are not cached
	fake.err = errors.New("toggl down")
	_, _ = c.Projects(ctx)
	_, _ = c.Projects(ctx)
	if fake.calls != 6 {
		t.Errorf("error: %d backend calls, want 6", fake.calls)
	}
	fake.err = timenote.ErrNoCurrentEntry
	_ = c.Done(ctx) // clears the cache
	_, _ = c.Current(ctx)
	_, err := c.Current(ctx)
	var problem *remote.Problem
	if fake.calls != 8 || !errors.As(err, &problem) || problem.Code != "NO_CURRENT_ENTRY" {
		t.Errorf("no current entry: %d backend calls, %v", fake.calls, err)
	}
}

func TestErrors(t *testing.T) {
	ctx := t.Context()
	fake := &fakeBackend{}
	ts := httptest.NewServer(newHandler(fake, "secret", 0))
	defer ts.Close()

	var problem *remote.Problem

	err := remote.New(ts.URL, "wrong").New(ctx)
	if !errors.As(err, &problem) || problem.Status != http.StatusUnauthorized || problem.Code != "UNAUTHORIZED" {
		t.Errorf("wrong token: %v", err)
	}
	if fake.last != "" {
		t.Errorf("backend called without auth: %q", fake.last)
	}

	c := remote.New(ts.URL, "secret")
	err = c.CreateProject(ctx, "")
	if !errors.As(err, &problem) || problem.Status != http.StatusUnprocessableEntity {
		t.Errorf("empty name: %v", err)
	}

	fake.err = fmt.Errorf("wrapped: %w", timenote.ErrNoCurrentEntry)
	err = c.Done(ctx)
	if !errors.As(err, &problem) || problem.Status != http.StatusNotFound || problem.Code != "NO_CURRENT_ENTRY" {
		t.Errorf("no current entry: %v", err)
	}

	fake.err = errors.New("toggl down")
	err = c.Done(ctx)
	if !errors.As(err, &problem) || problem.Status != http.StatusBadGateway || problem.Detail != "toggl down" {
		t.Errorf("backend error: %v", err)
	}
	if _, err := c.Current(ctx); err == nil {
		t.Error("Current: expected error")
	}
}
