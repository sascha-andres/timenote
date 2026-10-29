package cmd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/jason0x43/go-toggl"
)

func TestWizardNewWithDescription(t *testing.T) {
	c1, c2 := 1, 2
	clients := []toggl.Client{{ID: c1, Name: "first"}, {ID: c2, Name: "second"}}
	projects := []toggl.Project{
		{ID: 10, Name: "p-first", Cid: &c1},
		{ID: 20, Name: "p-second", Cid: &c2},
		{ID: 30, Name: "p-none"},
	}
	var m tea.Model = newWizardModel(nil, clients, projects)
	press := func(k tea.KeyPressMsg) {
		m, _ = m.Update(k)
	}
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	down := tea.KeyPressMsg{Code: tea.KeyDown}

	press(enter) // new
	if got := len(m.(wizardModel).list.Items()); got != 3 {
		t.Fatalf("client entries = %d, want 3 (two clients and no client)", got)
	}
	press(down)
	press(enter) // second client
	press(enter) // p-second
	press(enter) // add description: yes
	for _, r := range "work" {
		press(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	press(enter)

	w := m.(wizardModel)
	if !w.done || w.action != actionNew || w.project != "p-second" || w.description != "work" {
		t.Errorf("got done=%v action=%q project=%q description=%q", w.done, w.action, w.project, w.description)
	}
}

func TestClientProjects(t *testing.T) {
	c1 := 1
	projects := []toggl.Project{{ID: 10, Cid: &c1}, {ID: 30}}
	if got := clientProjects(projects, 1); len(got) != 1 || got[0].ID != 10 {
		t.Errorf("clientProjects(1) = %v", got)
	}
	if got := clientProjects(projects, 0); len(got) != 1 || got[0].ID != 30 {
		t.Errorf("clientProjects(0) = %v", got)
	}
}
