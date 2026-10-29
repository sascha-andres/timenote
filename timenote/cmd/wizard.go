package cmd

import (
	"context"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"go.livingit.de/timenote"
)

type wizardStep int

const (
	stepAction wizardStep = iota
	stepChange
	stepClient
	stepProject
	stepAskDescription
	stepDescription
)

// Menu entries; the selected entry is stored as the wizard action.
const (
	actionNew       = "new"
	actionChange    = "change"
	actionStop      = "stop"
	actionAppend    = "append description"
	actionOverwrite = "overwrite description"
	actionProject   = "new project"
)

type (
	// choice is a list entry; id holds the backend id for clients.
	choice struct {
		title string
		id    int
	}

	// runningMsg reports the result of the running time entry check.
	runningMsg struct{ err error }

	wizardModel struct {
		p           timenote.Backend
		clients     []timenote.Client
		projects    []timenote.Project
		list        list.Model
		input       textinput.Model
		step        wizardStep
		action      string
		project     string
		description string
		done        bool
		err         error
	}
)

func (c choice) FilterValue() string { return c.title }
func (c choice) Title() string       { return c.title }
func (c choice) Description() string { return "" }

// runWizard asks step by step what to do and then starts, changes or stops
// the running time entry.
func runWizard() {
	p := newPersistor()
	clients, err := p.Clients(context.Background())
	if err != nil {
		fatal(err)
	}
	projects, err := p.Projects(context.Background())
	if err != nil {
		fatal(err)
	}
	res, err := tea.NewProgram(newWizardModel(p, clients, filterProjects(projects))).Run()
	if err != nil {
		fatal(err)
	}
	m := res.(wizardModel)
	if m.err != nil {
		fatal(m.err)
	}
	if !m.done {
		return
	}
	if err := m.apply(); err != nil {
		fatal(err)
	}
}

func newWizardModel(p timenote.Backend, clients []timenote.Client, projects []timenote.Project) wizardModel {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	m := wizardModel{
		p:        p,
		clients:  clients,
		projects: projects,
		list:     list.New(nil, d, 80, 20),
		input:    textinput.New(),
	}
	m.showMenu(stepAction, "timenote", actionNew, actionChange, actionStop)
	return m
}

// apply runs the action selected in the wizard.
func (m wizardModel) apply() error {
	switch m.action {
	case actionNew:
		if err := m.p.New(context.Background()); err != nil {
			return err
		}
		// Append writes the full cached entry, so it must run before the project is set.
		if m.description != "" {
			if err := m.p.Append(context.Background(), m.description, *separator); err != nil {
				return err
			}
		}
		return m.p.SetProjectForCurrentTimestamp(context.Background(), m.project, false)
	case actionAppend:
		return m.p.Append(context.Background(), m.description, *separator)
	case actionOverwrite:
		return m.p.SetDescription(context.Background(), m.description)
	case actionProject:
		return m.p.SetProjectForCurrentTimestamp(context.Background(), m.project, false)
	case actionStop:
		return m.p.Done(context.Background())
	}
	return nil
}

func (m wizardModel) Init() tea.Cmd {
	return nil
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case runningMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		return m, m.showMenu(stepChange, actionChange, actionAppend, actionOverwrite, actionProject)
	case tea.KeyPressMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "enter" && m.step == stepDescription:
			if m.description = m.input.Value(); m.description == "" {
				return m, nil
			}
			return m.finish()
		case msg.String() == "enter" && !m.list.SettingFilter():
			if c, ok := m.list.SelectedItem().(choice); ok {
				return m.choose(c)
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	if m.step == stepDescription {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m wizardModel) View() tea.View {
	content := m.list.View()
	if m.step == stepDescription {
		content = "description:\n\n" + m.input.View() + "\n\n(enter: confirm, ctrl+c: cancel)"
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// choose moves to the next step for the selected list entry.
func (m wizardModel) choose(c choice) (tea.Model, tea.Cmd) {
	switch m.step {
	case stepAction:
		if c.title == actionChange {
			return m, m.checkRunning
		}
		m.action = c.title
		if m.action == actionStop {
			return m.finish()
		}
		return m, m.showClients()
	case stepChange:
		m.action = c.title
		if m.action == actionProject {
			return m, m.showClients()
		}
		return m, m.showInput()
	case stepClient:
		return m, m.showProjects(c.id)
	case stepProject:
		m.project = c.title
		if m.action == actionNew {
			return m, m.showMenu(stepAskDescription, "add description?", "yes", "no")
		}
		return m.finish()
	case stepAskDescription:
		if c.title == "yes" {
			return m, m.showInput()
		}
		return m.finish()
	}
	return m, nil
}

func (m wizardModel) finish() (tea.Model, tea.Cmd) {
	m.done = true
	return m, tea.Quit
}

// checkRunning reads the running time entry from the backend, not from the cache.
func (m wizardModel) checkRunning() tea.Msg {
	_, err := m.p.Current(context.Background())
	return runningMsg{err: err}
}

func (m *wizardModel) showMenu(step wizardStep, title string, entries ...string) tea.Cmd {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = choice{title: e}
	}
	return m.showList(step, title, items)
}

func (m *wizardModel) showClients() tea.Cmd {
	items := make([]list.Item, 0, len(m.clients)+1)
	for _, c := range m.clients {
		items = append(items, choice{title: c.Name, id: c.ID})
	}
	if len(clientProjects(m.projects, 0)) > 0 {
		items = append(items, choice{title: "(no client)"})
	}
	return m.showList(stepClient, "client", items)
}

func (m *wizardModel) showProjects(clientID int) tea.Cmd {
	projects := clientProjects(m.projects, clientID)
	items := make([]list.Item, len(projects))
	for i, prj := range projects {
		items[i] = choice{title: prj.Name, id: prj.ID}
	}
	return m.showList(stepProject, "project", items)
}

func (m *wizardModel) showList(step wizardStep, title string, items []list.Item) tea.Cmd {
	m.step = step
	m.list.Title = title
	m.list.ResetFilter()
	m.list.ResetSelected()
	return m.list.SetItems(items)
}

func (m *wizardModel) showInput() tea.Cmd {
	m.step = stepDescription
	m.input.Reset()
	return m.input.Focus()
}

// clientProjects returns the projects of clientID; 0 selects projects without a client.
func clientProjects(projects []timenote.Project, clientID int) []timenote.Project {
	var res []timenote.Project
	for _, prj := range projects {
		if prj.ClientID == clientID {
			res = append(res, prj)
		}
	}
	return res
}
