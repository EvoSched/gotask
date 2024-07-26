package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case Main, Pending, Archived:
			return m.updateTable(msg, msg.String())
		case Detail:
			m.updateDetail(msg.String())
		case List:
			return m.updateFile(msg, msg.String())
		case Help:
			return m.updateHelp(msg, msg.String())
		}
	}
	return m, cmd
}

func (m model) updateTable(msg tea.Msg, s string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch s {
	case "t":
		m.state = Main
		m.main, cmd = m.main.Update(msg)
	case "p":
		m.state = Pending
		m.pend, cmd = m.pend.Update(msg)
	case "f":
		m.state = Archived
		m.arch, cmd = m.arch.Update(msg)
	case "a":
		// add task
		//t := append(m.main.Rows(), nil)
		//m.main.SetRows(t)
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
		if m.state == Main {
			revSort(&m.main)
			m.main, cmd = m.main.Update(msg)
		} else if m.state == Pending {
			revSort(&m.pend)
			m.pend, cmd = m.pend.Update(msg)
		} else {
			revSort(&m.arch)
			m.arch, cmd = m.arch.Update(msg)
		}
	case "tab", "2":
		m.state = List
		m.list, cmd = m.list.Update(msg)
	case "h":
		m.state = Help
	case "esc":
		// focus or blur table
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) updateDetail(s string) {
	switch s {
	case "m":
	// modify task
	case "q":
		// return to task list
	}
}

func (m model) updateFile(msg tea.Msg, s string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch s {
	case "enter":
	// setup option for filename
	case "1":
		m.state = Main
		m.main, cmd = m.main.Update(msg)
	case "h":
		m.state = Help
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, cmd
}

func (m model) updateHelp(msg tea.Msg, s string) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch s {
	case "1", "tab":
		m.state = Main
		m.main, cmd = m.main.Update(msg)
	case "2":
		m.state = List
		m.list, cmd = m.list.Update(msg)
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, cmd
}

func revSort(table *table.Model) {
	for i := 0; i < len(table.Rows()); i++ {
		t := table.Rows()[i]
		table.Rows()[i] = table.Rows()[len(table.Rows())-i-1]
		table.Rows()[len(table.Rows())-i-1] = t
	}
}
