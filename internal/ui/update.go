package ui

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case Main, Pending, Archived:
			m.updateTable(msg.String())
		case Detail:
			m.updateDetail(msg.String())
		case List:
			m.updateFile(msg.String())
		case Help:
			m.updateHelp(msg.String())
		}
	}
	return nil, nil
}

func (m model) updateTable(s string) {
	switch s {
	case "t":
		m.state = Main
	case "p":
		m.state = Pending
	case "f":
		m.state = Archived
	case "a":
	// add task
	case "d":
	// marks task as done
	case "u":
	// marks task as not done
	case "m":
	// modify task
	case "r":
	// removes task
	case "n":
	// add note to task (opens up a text input area)
	case "enter":
	// views task details
	case "s":
	// sorts task table in ascending or descending order
	case "tab", "2":
		// switches over to file state
		m.state = List
	case "h":
	// sets state to help page
	case "ctrl+c", "q":
		// exits application
	}
}

func (m model) updateDetail(s string) {
	switch s {
	case "m":
	// modify task
	case "q":
		// return to task list
	}
}

func (m model) updateFile(s string) {
	switch s {
	case "enter":
	// setup option for filename
	case "1":
	case "2":
	case "h":
	case "ctrl+c", "q":
	}
}

func (m model) updateHelp(s string) {
	switch s {
	case "1", "tab":
		m.state = Main
	case "2":
		m.state = List
	case "ctrl+c", "q":
		// exits application
	}
}
